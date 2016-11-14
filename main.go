package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

var (
	Version string
	Build   string
)

var (
	version  = flag.Bool("version", false, "Version of docker-volume-location-plugin")
	location = flag.String("location", "/mnt/docker-volumes", "Where to save docker-volumes")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("docker-volume-location-plugin version: %s, build: %s\n", Version, Build)
		os.Exit(0)
	}

	d := newVolumeLocationDriver(location)
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "vlp"))
}
