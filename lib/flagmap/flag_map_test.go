package flagmap_test

import (
	"testing"

	"github.com/npetzall/docker-volume-location-plugin/lib/flagmap"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFlagMapSetAliasAndPath(t *testing.T) {
	Convey("Create a flagmap", t, func() {
		var f = make(flagmap.FlagMap)
		Convey("Given input string \"alias=path\"", func() {
			i := "alias=path"
			Convey("When added through set", func() {
				f.Set(i)
				Convey("There should be a key \"alias\" with value \"path\"", func() {
					v, ke := f["alias"]
					So(ke, ShouldBeTrue)
					So(v, ShouldEqual, "path")
				})
			})
		})
		Convey("Given input string \"path\"", func() {
			i := "path"
			Convey("When added through set", func() {
				f.Set(i)
				Convey("There should be a key \"default\" with value \"path\"", func() {
					v, ke := f["default"]
					So(ke, ShouldEqual, true)
					So(v, ShouldEqual, "path")
				})
			})
		})
		Convey("Given invalid input string \"alias=path=invalid\"", func() {
			i := "alias=path=invalid"
			Convey("When added through set", func() {
				f.Set(i)
				Convey("Then no key or values should exist", func() {
					So(f, ShouldBeEmpty)
				})
			})
		})
	})
}
