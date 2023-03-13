package main

type module struct { // defined modules for installation in the docker file
	gitrepo string   // git repo url
	usegit  bool     // clone the git repo to build/install?
	aptadds []string // dependencies to get through apt
	build   []string // actual build commands outside of cloning and moving into the repo
}

func modules() map[string]module {
	mods := make(map[string]module)
	mods["rizin"] = module{
		gitrepo: "https://github.com/rizinorg/rizin",
		usegit:  true,
		aptadds: []string{"meson", "ninja-build"},
		build:   []string{"meson --buildtype=release build", "ninja -C build", "ninja -C build install"},
	}
	return mods
}

func main() {
	image := "ubuntu:jammy"
	mods := modules()

}
