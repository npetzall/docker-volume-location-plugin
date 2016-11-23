package driver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/npetzall/docker-volume-location-plugin/lib/profiler"
)

type volumeDesc struct {
	alias    *string
	location *string
	name     string
}

func (vd *volumeDesc) mountpoint() string {
	return filepath.Join(*vd.location, vd.name)
}

type volumeLocationDriver struct {
	locations *map[string]string
	volumes   map[string]*volumeDesc
	m         *sync.Mutex
}

//Create a new volumeLocationDriver
//Supply a *map[string]string where key is alias and value is a path
func NewVolumeLocationDriver(locations *map[string]string) volumeLocationDriver {
	d := volumeLocationDriver{
		locations: locations,
		volumes:   map[string]*volumeDesc{},
		m:         &sync.Mutex{},
	}
	d.m.Lock()
	defer d.m.Unlock()
	d.load()
	return d
}

func (d volumeLocationDriver) load() {
	for alias, location := range *d.locations {
		d.loadLocation(alias, location)
	}
	return
}

func (d volumeLocationDriver) loadLocation(alias string, location string) {
	defer profiler.Timed(time.Now(), "load")
	files, err := ioutil.ReadDir(location)
	if err != nil {
		fmt.Printf("Unable to read location: %s\n", location)
		return
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			if _, exist := d.volumes[fileInfo.Name()]; !exist {
				d.volumes[fileInfo.Name()] = &volumeDesc{
					alias:    &alias,
					location: &location,
					name:     fileInfo.Name(),
				}
			}
		}
	}
	return
}

func (d volumeLocationDriver) Create(r volume.Request) volume.Response {
	defer profiler.Timed(time.Now(), "Create")
	d.m.Lock()
	defer d.m.Unlock()

	var alias string
	var location string

	if ov, exist := r.Options["location"]; exist {
		alias = ov
	} else {
		alias = "default"
	}

	if loc, exist := (*d.locations)[alias]; exist {
		location = loc
	} else {
		return volume.Response{Err: fmt.Sprintf("Unable to find location for \"%s\"", alias)}
	}

	// Exists do nothing
	if vol, exist := d.volumes[r.Name]; exist {
		if *vol.location != location {
			return volume.Response{Err: "Volume exists but another location"}
		}
		return volume.Response{}
	}

	//Create
	err := os.MkdirAll(filepath.Join(location, r.Name), 0755)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	//Add to driver volumes
	d.volumes[r.Name] = &volumeDesc{
		alias:    &alias,
		location: &location,
		name:     r.Name,
	}
	//log.Printf("Created volume %s\n", r.Name)
	return volume.Response{}
}

func (d volumeLocationDriver) List(r volume.Request) volume.Response {
	defer profiler.Timed(time.Now(), "List")
	d.m.Lock()
	defer d.m.Unlock()
	res := volume.Response{}
	d.load()
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
	defer profiler.Timed(time.Now(), "Get")
	d.m.Lock()
	defer d.m.Unlock()
	if vol, exist := d.volumes[r.Name]; exist {
		return volume.Response{Volume: &volume.Volume{Name: vol.name, Mountpoint: vol.mountpoint()}}
	}
	return volume.Response{Err: fmt.Sprintf("Unable to find volume: %s", r.Name)}
}

func (d volumeLocationDriver) Remove(r volume.Request) volume.Response {
	defer profiler.Timed(time.Now(), "Remove")
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
	defer profiler.Timed(time.Now(), "Path")
	d.m.Lock()
	defer d.m.Unlock()
	return d.mountPoint(r.Name)
}

func (d volumeLocationDriver) mountPoint(name string) volume.Response {
	defer profiler.Timed(time.Now(), "MountPoint")
	if vol, exist := d.volumes[name]; exist {
		return volume.Response{Mountpoint: vol.mountpoint()}
	}
	return volume.Response{Err: fmt.Sprintf("Unable to find volume: %s", name)}
}

func (d volumeLocationDriver) Mount(r volume.MountRequest) volume.Response {
	defer profiler.Timed(time.Now(), "Mount")
	d.m.Lock()
	defer d.m.Unlock()
	return d.mountPoint(r.Name)
}

func (d volumeLocationDriver) Unmount(r volume.UnmountRequest) volume.Response {
	defer profiler.Timed(time.Now(), "Unmount")
	d.m.Lock()
	defer d.m.Unlock()
	return d.mountPoint(r.Name)
}

func (d volumeLocationDriver) Capabilities(r volume.Request) volume.Response {
	defer profiler.Timed(time.Now(), "Capabilities")
	var res volume.Response
	res.Capabilities = volume.Capability{Scope: "global"}
	return res
}
