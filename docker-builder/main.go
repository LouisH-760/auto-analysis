package main

import (
	"flag"
	"fmt"
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
	mods["rizin"] = module{
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
		build: []string{
			"wget https://www.clamav.net/downloads/production/clamav-1.0.1.linux.x86_64.deb -O clamav.deb",
			"sudo dpkg -i clamav.deb",
			"rm clamav.deb",
		},
	}
	return mods
}

func aptList(mods map[string]module) []string {
	pkgs := []string{"git", "python3", "golang-go"}
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
		fcmd = append(fcmd, fmt.Sprintf("RUN git clone %s %s\nWORKDIR %s", mod.gitrepo, name, name))
	}
	for _, cmd := range mod.build {
		fcmd = append(fcmd, fmt.Sprintf("RUN %s", cmd))
	}
	if mod.usegit {
		fcmd = append(fcmd, "WORKDIR ..")
		fcmd = append(fcmd, fmt.Sprintf("RUN rm -rf %s", name))
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
	// apt stuff
	pkgs := aptList(mods)
	aptcmd := fmt.Sprintf("RUN apt update && apt install -y %s", strings.Join(pkgs, " "))
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
	df := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", header, workdir, copysample, aptcmd, icmd, fbuild, expose, entry)
	return df
}

func main() {
	image := "ubuntu:jammy"
	implantrepo := "https://github.com/LouisH-760/auto-analysis"
	mods := modules()
	sample := flag.String("sample", "", "Path to the sample")
	modfolder := flag.String("modules", "./modules", "path to the modules folder. Default: ./modules")
	port := 8080
	for name, mod := range mods {
		// can't do mod[name].used directly, so use the object copy and assign it in place
		mod.used = flag.Bool(name, false, fmt.Sprintf("Build the docker container with support for the %s module (%s)", name, mod.gitrepo))
		mods[name] = mod
	}
	flag.Parse()
	if len(*sample) < 1 {
		fmt.Print("The -sample flag is mandatory! and may not be empty\n")
		return
	}
	dockerfile := dockerFile(image, *sample, *modfolder, mods, port, implantrepo)
	fmt.Printf("%s\n", dockerfile)
}
