package driver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/npetzall/docker-volume-location-plugin/lib/driver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDriver(t *testing.T) {
	basePaht := filepath.Join(os.TempDir(), "vlp-test")
	os.RemoveAll(basePaht)
	locations := map[string]string{
		"default": filepath.Join(basePaht, "dvdef"),
		"dv1":     filepath.Join(basePaht, "dv1"),
		"dv2":     filepath.Join(basePaht, "dv2"),
		"dv3":     filepath.Join(basePaht, "dv3"),
		"dv4":     filepath.Join(basePaht, "dv4"),
	}
	for _, location := range locations {
		os.MkdirAll(location, 0755)
	}
	os.MkdirAll(filepath.Join(locations["default"], "exists"), 0755)
	d := driver.NewVolumeLocationDriver(&locations)
	Convey("Setup", t, func() {
		Convey("Existing volumes are found", func() {
			res := d.List(volume.Request{})
			So(len(res.Volumes), ShouldEqual, 1)
			So(res.Volumes[0].Name, ShouldEqual, "exists")
			So(res.Volumes[0].Mountpoint, ShouldEqual, filepath.Join(locations["default"], "exists"))
		})
		Convey("Create new volume at default", func() {
			res := d.Create(volume.Request{Name: "atDefault"})
			So(res.Err, ShouldBeEmpty)
			fileInfo, err := os.Stat(filepath.Join(locations["default"], "atDefault"))
			So(err, ShouldBeNil)
			So(fileInfo.IsDir(), ShouldBeTrue)
			Convey("Volume shows up in List", func() {
				res := d.List(volume.Request{})
				So(res.Volumes, ShouldContain, &volume.Volume{Name: "atDefault", Mountpoint: filepath.Join(locations["default"], "atDefault")})
				Convey("Remove volume", func() {
					res := d.Remove(volume.Request{Name: "atDefault"})
					So(res.Err, ShouldBeEmpty)
					Convey("Doesn't show up in List anymore", func() {
						res := d.List(volume.Request{})
						So(res.Volumes, ShouldNotContain, &volume.Volume{Name: "atDefault", Mountpoint: filepath.Join(locations["default"], "atDefault")})
					})
				})
			})
		})
		Convey("Create volume at dv1", func() {
			res := d.Create(volume.Request{Name: "atDv1", Options: map[string]string{"location": "dv1"}})
			So(res.Err, ShouldBeEmpty)
			fileInfo, err := os.Stat(filepath.Join(locations["dv1"], "atDv1"))
			So(err, ShouldBeNil)
			So(fileInfo.IsDir(), ShouldBeTrue)
			Convey("Volume shows up in List", func() {
				res := d.List(volume.Request{})
				So(res.Volumes, ShouldContain, &volume.Volume{Name: "atDv1", Mountpoint: filepath.Join(locations["dv1"], "atDv1")})
				Convey("Remove volume", func() {
					res := d.Remove(volume.Request{Name: "atDv1"})
					So(res.Err, ShouldBeEmpty)
					Convey("Doesn't show up in List anymore", func() {
						res := d.List(volume.Request{})
						So(res.Volumes, ShouldNotContain, &volume.Volume{Name: "atDv1", Mountpoint: filepath.Join(locations["dv1"], "atDv1")})
					})
				})
			})
		})
		Convey("Create volume at dv2", func() {
			res := d.Create(volume.Request{Name: "atDv2", Options: map[string]string{"location": "dv2"}})
			So(res.Err, ShouldBeEmpty)
			fileInfo, err := os.Stat(filepath.Join(locations["dv2"], "atDv2"))
			So(err, ShouldBeNil)
			So(fileInfo.IsDir(), ShouldBeTrue)
			Convey("Volume shows up in List", func() {
				res := d.List(volume.Request{})
				So(res.Err, ShouldBeEmpty)
				Convey("Create at same locations returns no error (atDv2)", func() {
					res := d.Create(volume.Request{Name: "atDv2", Options: map[string]string{"location": "dv2"}})
					So(res.Err, ShouldBeEmpty)
					Convey("Create at other location returns error (atDv2 at dv1)", func() {
						res := d.Create(volume.Request{Name: "atDv2", Options: map[string]string{"location": "dv1"}})
						So(res.Err, ShouldNotBeEmpty)
					})
				})
			})
		})
		Convey("Create volume at dv3", func() {
			res := d.Create(volume.Request{Name: "atDv3", Options: map[string]string{"location": "dv3"}})
			So(res.Err, ShouldBeEmpty)
			fileInfo, err := os.Stat(filepath.Join(locations["dv3"], "atDv3"))
			So(err, ShouldBeNil)
			So(fileInfo.IsDir(), ShouldBeTrue)
			Convey("Volume shows up in List", func() {
				res := d.List(volume.Request{})
				So(res.Err, ShouldBeEmpty)
				Convey("Get of volume will return correct mountpoint", func() {
					res := d.Get(volume.Request{Name: "atDv3"})
					So(res.Err, ShouldBeEmpty)
				})
			})
		})
		Convey("Create volume at dv4", func() {
			res := d.Create(volume.Request{Name: "atDv4", Options: map[string]string{"location": "dv4"}})
			So(res.Err, ShouldBeEmpty)
			fileInfo, err := os.Stat(filepath.Join(locations["dv4"], "atDv4"))
			So(err, ShouldBeNil)
			So(fileInfo.IsDir(), ShouldBeTrue)
			Convey("Get path and verify path", func() {
				res := d.Path(volume.Request{Name: "atDv4"})
				So(res.Err, ShouldBeEmpty)
				So(res.Mountpoint, ShouldEqual, filepath.Join(locations["dv4"], "atDv4"))
			})
		})
		Convey("Mount return mountpoint", func() {
			res := d.Mount(volume.MountRequest{Name: "exists", ID: "1"})
			So(res.Err, ShouldBeEmpty)
			Convey("Unmount", func() {
				res := d.Unmount(volume.UnmountRequest{Name: "exists", ID: "1"})
				So(res.Err, ShouldBeEmpty)
			})
		})
		Convey("Get for unknown volume will return error", func() {
			res := d.Get(volume.Request{Name: "atDv3DoesntExist"})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Create volume with alias that doesn't exist causes error", func() {
			res := d.Create(volume.Request{Name: "shouldntBeCreated", Options: map[string]string{"location": "doesntExist"}})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Remove a volume that doesn't exist, returns error", func() {
			res := d.Remove(volume.Request{Name: "doesntExist"})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Path for non existent volume, returns error", func() {
			res := d.Path(volume.Request{Name: "doesntExist"})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Mount non existet", func() {
			res := d.Mount(volume.MountRequest{Name: "doesntExist", ID: "2"})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Unmount non existet", func() {
			res := d.Unmount(volume.UnmountRequest{Name: "doesntExist", ID: "2"})
			So(res.Err, ShouldNotBeEmpty)
		})
		Convey("Capabilities should contain Scope:Global", func() {
			res := d.Capabilities(volume.Request{})
			So(res.Err, ShouldBeEmpty)
			So(res.Capabilities.Scope, ShouldEqual, "global")
		})
	})
	os.RemoveAll(basePaht)
}
