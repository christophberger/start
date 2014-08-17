package start

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParse(t *testing.T) {
	SkipConvey("", t, func() {
		var x int = 1
		Convey("", func() {
			x++

			Convey("", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}

func TestUp(t *testing.T) {
	SkipConvey("", t, func() {
		var x int = 1
		Convey("", func() {
			x++

			Convey("", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
