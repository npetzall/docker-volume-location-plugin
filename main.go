package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/npetzall/docker-volume-location-plugin/lib/driver"
	"github.com/npetzall/docker-volume-location-plugin/lib/flagmap"
	"github.com/npetzall/docker-volume-location-plugin/lib/profiler"
)

var (
	//Version set by build tool
	Version string
	//Build revision set by build tool
	Build string
)

var (
	locations = make(flagmap.FlagMap)
	profile   = flag.Bool("profile", false, "profile executions")
	version   = flag.Bool("version", false, "Version of docker-volume-location-plugin")
)

func init() {
	flag.Var(&locations, "location", "[[alias=]path] can be declared multiple times\n\tomitting alias= sets default")
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-location [alias=]/mnt/docker-volumes]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *profile {
		profiler.SetProfiling(*profile)
	}

	if *version {
		fmt.Printf("docker-volume-location-plugin version: %s, build: %s\n", Version, Build)
		os.Exit(0)
	}

	d := driver.NewVolumeLocationDriver((*map[string]string)(&locations))
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "vlp"))
}
