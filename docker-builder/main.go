package main

import (
	"flag"
	"fmt"
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
		build:   []string{"meson --buildtype=release --prefix=/usr build", "ninja -C build", "ninja -C build install"},
	}
	mods["volatility2"] = module{ // specifically get the last volatility 2 standalone version
		gitrepo: "https://github.com/volatilityfoundation/volatility.git",
		usegit:  false,
		aptadds: []string{"unzip"},
		build:   []string{"wget http://downloads.volatilityfoundation.org/releases/2.6/volatility_2.6_lin64_standalone.zip -O volatility.zip", "unzip volatility.zip", "mv ./volatility_2.6_lin64_standalone/volatility_2.6_lin64_standalone /usr/bin/volatility2", "rm volatility.zip"},
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
		build:   []string{"wget https://www.clamav.net/downloads/production/clamav-1.0.1.linux.x86_64.deb -O clamav.deb", "sudo dpkg -i clamav.deb", "rm clamav.deb"},
	}
	return mods
}

func main() {
	image := "ubuntu:jammy"
	mods := modules()
	sample := flag.String("sample", "", "Path to the sample")
	modfolder := flag.String("modules", "./modules", "path to the modules folder. Default: ./modules")
	for name, mod := range mods {
		mod.used = flag.Bool(name, false, fmt.Sprintf("Build the docker container with support for the %s module (%s)", name, mod.gitrepo)) // can't do mod[name].used directly, so use the object copy and assign it in place
		mods[name] = mod
	}
	flag.Parse()
}
