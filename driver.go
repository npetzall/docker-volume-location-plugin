package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
)

func timed(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

type volumeDesc struct {
	location    *string
	name        string
	connections int
}

func (vd *volumeDesc) mountpoint() string {
	return filepath.Join(*vd.location, vd.name)
}

type volumeLocationDriver struct {
	location *string
	volumes  map[string]*volumeDesc
	m        *sync.Mutex
}

func newVolumeLocationDriver(location *string) volumeLocationDriver {
	d := volumeLocationDriver{
		location: location,
		volumes:  map[string]*volumeDesc{},
		m:        &sync.Mutex{},
	}
	d.m.Lock()
	defer d.m.Unlock()
	err := d.load()
	if err != nil {
		fmt.Println(err)
	}
	return d
}

func (d volumeLocationDriver) load() error {
	defer timed(time.Now(), "load")
	files, err := ioutil.ReadDir(*d.location)
	if err != nil {
		fmt.Printf("Unable to read location: %s\n", *d.location)
		return err
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			if _, exist := d.volumes[fileInfo.Name()]; !exist {
				d.volumes[fileInfo.Name()] = &volumeDesc{
					location: d.location,
					name:     fileInfo.Name(),
				}
			}
		}
	}
	return nil
}

func (d volumeLocationDriver) Create(r volume.Request) volume.Response {
	defer timed(time.Now(), "Create")
	log.Printf("Creating volume %s\n", r.Name)
	d.m.Lock()
	defer d.m.Unlock()

	// Exists do nothing
	if _, exist := d.volumes[r.Name]; exist {
		return volume.Response{}
	}

	//Create
	err := os.MkdirAll(filepath.Join(*d.location, r.Name), 0755)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	//Add to driver volumes
	d.volumes[r.Name] = &volumeDesc{
		location: d.location,
		name:     r.Name,
	}
	return volume.Response{}
}

func (d volumeLocationDriver) List(r volume.Request) volume.Response {
	defer timed(time.Now(), "List")
	d.m.Lock()
	defer d.m.Unlock()
	res := volume.Response{}
	if err := d.load(); err != nil {
		res.Err = err.Error()
	}
	res.Volumes = make([]*volume.Volume, len(d.volumes))
	var index = 0
	for _, vol := range d.volumes {
		res.Volumes[index] = &volume.Volume{
			Name:       vol.name,
			Mountpoint: vol.mountpoint(),
		}
		index++
	}
	return res
}

func (d volumeLocationDriver) Get(r volume.Request) volume.Response {
	defer timed(time.Now(), "Get")
	d.m.Lock()
	defer d.m.Unlock()
	if vol, exist := d.volumes[r.Name]; exist {
		return volume.Response{Volume: &volume.Volume{Name: vol.name, Mountpoint: vol.mountpoint()}}
	}
	return volume.Response{Err: fmt.Sprintf("Unable to find volume: %s", r.Name)}
}

func (d volumeLocationDriver) Remove(r volume.Request) volume.Response {
	defer timed(time.Now(), "Remove")
	d.m.Lock()
	defer d.m.Unlock()
	if vol, exist := d.volumes[r.Name]; exist {
		os.RemoveAll(vol.mountpoint())
		delete(d.volumes, r.Name)
		return volume.Response{}
	}
	return volume.Response{Err: fmt.Sprintf("Unable to find volume: %s", r.Name)}
}

func (d volumeLocationDriver) Path(r volume.Request) volume.Response {
	defer timed(time.Now(), "Path")
	d.m.Lock()
	defer d.m.Unlock()
	return d.mountPoint(r.Name)
}

func (d volumeLocationDriver) mountPoint(name string) volume.Response {
	defer timed(time.Now(), "MountPoint")
	if vol, exist := d.volumes[name]; exist {
		return volume.Response{Mountpoint: vol.mountpoint()}
	}
	return volume.Response{Err: fmt.Sprintf("Unable to find volume: %s", name)}
}

func (d volumeLocationDriver) Mount(r volume.MountRequest) volume.Response {
	defer timed(time.Now(), "Mount")
	d.m.Lock()
	defer d.m.Unlock()
	res := d.mountPoint(r.Name)
	if res.Err == "" {
		d.volumes[r.Name].connections++
	}
	return res
}

func (d volumeLocationDriver) Unmount(r volume.UnmountRequest) volume.Response {
	defer timed(time.Now(), "Unmount")
	d.m.Lock()
	defer d.m.Unlock()
	res := d.mountPoint(r.Name)
	if res.Err == "" {
		d.volumes[r.Name].connections--
	}
	return res
}

func (d volumeLocationDriver) Capabilities(r volume.Request) volume.Response {
	defer timed(time.Now(), "Capabilities")
	var res volume.Response
	res.Capabilities = volume.Capability{Scope: "global"}
	return res
}
