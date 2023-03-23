package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type module struct { // defined modules for installation in the docker file
	gitrepo string   // git repo url
	usegit  bool     // clone the git repo to build/install?
	aptadds []string // dependencies to get through apt
	build   []string // actual build commands outside of cloning and moving into the repo
	used    *bool    // wether the flag should be used, defined using arguments
}

func modules() map[string]module {
	mods := make(map[string]module)
	mods["rizin"] = module{ // should grab a specific version (deb?) instead of building from source from main
		gitrepo: "https://github.com/rizinorg/rizin",
		usegit:  true,
		aptadds: []string{"meson", "ninja-build"},
		build: []string{
			"meson --buildtype=release --prefix=/usr build",
			"ninja -C build",
			"ninja -C build install",
		},
	}
	mods["volatility2"] = module{ // specifically get the last volatility 2 standalone version
		gitrepo: "https://github.com/volatilityfoundation/volatility.git",
		usegit:  false,
		aptadds: []string{"unzip"},
		build: []string{
			"wget http://downloads.volatilityfoundation.org/releases/2.6/volatility_2.6_lin64_standalone.zip -O volatility.zip",
			"unzip volatility.zip",
			"mv ./volatility_2.6_lin64_standalone/volatility_2.6_lin64_standalone /usr/bin/volatility2",
			"rm volatility.zip",
		},
	}
	mods["yara"] = module{
		gitrepo: "https://github.com/VirusTotal/yara",
		usegit:  false,
		aptadds: []string{"yara"},
		build:   []string{},
	}
	mods["clamav"] = module{
		gitrepo: "https://github.com/Cisco-Talos/clamav",
		usegit:  false,
		aptadds: []string{},
		build: []string{ // grab a specific clamav version (latest at the time of writing)
			"wget https://www.clamav.net/downloads/production/clamav-1.0.1.linux.x86_64.deb -O clamav.deb",
			"dpkg -i clamav.deb",
			"rm clamav.deb",
		},
	}
	mods["diec"] = module{
		gitrepo: "https://github.com/horsicq/DIE-engine",
		usegit:  false,
		aptadds: []string{"libqt5script5", "libqt5scripttools5", "libqt5opengl5", "libqt5sql5"},
		build: []string{
			"wget https://github.com/horsicq/DIE-engine/releases/download/3.07/die_3.07_Ubuntu_22.04_amd64.deb -O die.deb",
			"dpkg -i die.deb",
			"rm die.deb",
		},
	}
	return mods
}

func aptList(mods map[string]module) []string {
	pkgs := []string{"git", "python3", "golang-go", "wget"}
	for _, mod := range mods {
		if *mod.used {
			pkgs = append(pkgs, mod.aptadds...)
		}
	}
	return pkgs
}

func buildcmd(mod module, name string) string {
	fcmd := []string{}
	if mod.usegit {
		fcmd = append(fcmd, fmt.Sprintf("RUN git clone --recursive \"%s\" \"%s\"\nWORKDIR \"%s\"", mod.gitrepo, name, name))
	}
	for _, cmd := range mod.build {
		fcmd = append(fcmd, fmt.Sprintf("RUN %s", cmd))
	}
	if mod.usegit {
		fcmd = append(fcmd, "WORKDIR ..")
		fcmd = append(fcmd, fmt.Sprintf("RUN rm -rf \"%s\"", name))
	}
	return strings.Join(fcmd, "\n")
}

func implantcmd(repo string) string {
	fcmd := []string{
		fmt.Sprintf("RUN git clone %s implant", repo),
		"WORKDIR implant/implant",
		"RUN go build .",
		"RUN mv implant /usr/bin",
		"WORKDIR ../..",
		"RUN rm implant -r",
	}
	return strings.Join(fcmd, "\n")
}

func dockerFile(image string, sample string, modfolder string, mods map[string]module, port int, implantrepo string) string {
	// image thing
	header := fmt.Sprintf("FROM %s AS auto-analysis", image)
	// file things
	workdir := "WORKDIR /autoa"
	copysample := fmt.Sprintf("COPY \"%s\" \"/autoa/%s\"", sample, path.Base(sample))
	copymods := fmt.Sprintf("COPY \"%s\" \"/autoa/modules/\"", modfolder)
	// apt stuff
	pkgs := aptList(mods)
	aptcmd := fmt.Sprintf("ARG DEBIAN_FRONTEND=noninteractive\nENV TZ=Europe/Paris\nRUN apt update && apt install -y %s", strings.Join(pkgs, " "))
	// build the implant
	icmd := implantcmd(implantrepo)
	// build commands for relevant packages
	buildcmds := []string{}
	for name, mod := range mods {
		if *mod.used && len(mod.build) > 0 {
			buildcmds = append(buildcmds, buildcmd(mod, name))
		}
	}
	fbuild := strings.Join(buildcmds, "\n")
	// port
	expose := fmt.Sprintf("EXPOSE %d", port)
	// entrypoint
	entry := "ENTRYPOINT [\"/usr/bin/implant\"]"
	// format the Dockerfile
	df := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", header, workdir, copysample, copymods, aptcmd, icmd, fbuild, expose, entry)
	return df
}

func runCmdStreamOutput(name string, args ...string) error {
	// https://stackoverflow.com/questions/30725751/streaming-commands-output-progress
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe() // get pipe from stdout
	if err != nil {
		fmt.Printf("Could not get stdout pipe: %s", err.Error())
		return err
	}
	cmd.Stderr = cmd.Stdout // set stderr output to show up on stdout
	cmd.Start()             // start the command
	pipeScan := bufio.NewScanner(stdout)
	for pipeScan.Scan() {
		line := pipeScan.Text()
		fmt.Println(line)
	}
	return cmd.Wait() // wait on the command to exit
}

func build() error {
	return runCmdStreamOutput("docker", "build", ".", "-t", "auto-analysis", "--no-cache")
}

func run(port int) error {
	return runCmdStreamOutput("docker", "run", "--rm", "--name", "auto-analysis", "-p", fmt.Sprintf("%d:%d", port, port), "auto-analysis")
}

func main() {
	image := "ubuntu:jammy"
	implantrepo := "https://github.com/LouisH-760/auto-analysis"
	mods := modules()
	sample := flag.String("sample", "", "Path to the sample. Must be in the docker context.")
	modfolder := flag.String("modules", "", "path to the modules folder. Must be in the docker context.")
	runct := flag.Bool("run", false, "Start the docker container")
	port := 8080
	for name, mod := range mods {
		// can't do mod[name].used directly, so use the object copy and assign it in place
		mod.used = flag.Bool(name, false, fmt.Sprintf("Build the docker container with support for the %s module (%s)", name, mod.gitrepo))
		mods[name] = mod
	}
	flag.Parse()
	if len(*sample) < 1 {
		fmt.Print("The -sample flag is mandatory and may not be empty\n")
		return
	}
	if len(*modfolder) < 1 {
		fmt.Print("The -modules flag is mandatory and may not be empty\n")
		return
	}
	// copy modules folder locally
	dockerfile := dockerFile(image, *sample, *modfolder, mods, port, implantrepo)
	if *runct {
		if err := os.WriteFile("Dockerfile", []byte(dockerfile), 0644); err != nil {
			fmt.Printf("Could not write Dockerfile: %s", err.Error())
			return
		}

		if err := build(); err != nil {
			fmt.Printf("Error while building the container: %s\n", err.Error())
			return
		}
		if err := run(port); err != nil {
			fmt.Printf("Error while running the container: %s\n", err.Error())
			return
		}
	} else {
		fmt.Printf("%s\n", dockerfile)
	}
}
