package envutil_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/beego-common/envutil"
)

func TestGet(t *testing.T) {
	Convey("Given Get", t, func() {
		Convey("When env var is set", func() {
			t.Setenv("TEST_KEY", "hello")
			So(envutil.Get("TEST_KEY", "default"), ShouldEqual, "hello")
		})

		Convey("When env var is not set", func() {
			So(envutil.Get("UNSET_KEY_XYZ", "fallback"), ShouldEqual, "fallback")
		})
	})
}

func TestMustGet(t *testing.T) {
	Convey("Given MustGet", t, func() {
		Convey("When env var is set", func() {
			t.Setenv("MUST_KEY", "value")
			So(envutil.MustGet("MUST_KEY"), ShouldEqual, "value")
		})

		Convey("When env var is absent it should panic", func() {
			So(func() { envutil.MustGet("ABSENT_KEY_XYZ") }, ShouldPanic)
		})
	})
}

func TestGetInt(t *testing.T) {
	Convey("Given GetInt", t, func() {
		Convey("When env var is a valid int", func() {
			t.Setenv("INT_KEY", "42")
			So(envutil.GetInt("INT_KEY", 0), ShouldEqual, 42)
		})

		Convey("When env var is not set", func() {
			So(envutil.GetInt("UNSET_INT_KEY", 7), ShouldEqual, 7)
		})

		Convey("When env var is not a valid int", func() {
			t.Setenv("BAD_INT", "not-a-number")
			So(envutil.GetInt("BAD_INT", 99), ShouldEqual, 99)
		})
	})
}

func TestGetBool(t *testing.T) {
	Convey("Given GetBool", t, func() {
		for _, truthy := range []string{"1", "true", "yes", "on", "TRUE", "YES"} {
			truthy := truthy
			Convey("When env var is truthy: "+truthy, func() {
				t.Setenv("BOOL_KEY", truthy)
				So(envutil.GetBool("BOOL_KEY", false), ShouldBeTrue)
			})
		}

		for _, falsy := range []string{"0", "false", "no", "off"} {
			falsy := falsy
			Convey("When env var is falsy: "+falsy, func() {
				t.Setenv("BOOL_KEY", falsy)
				So(envutil.GetBool("BOOL_KEY", true), ShouldBeFalse)
			})
		}

		Convey("When env var is not set", func() {
			So(envutil.GetBool("UNSET_BOOL_KEY", true), ShouldBeTrue)
		})
	})
}

func TestGetDuration(t *testing.T) {
	Convey("Given GetDuration", t, func() {
		Convey("When env var is a valid duration", func() {
			t.Setenv("DUR_KEY", "30s")
			So(envutil.GetDuration("DUR_KEY", time.Minute), ShouldEqual, 30*time.Second)
		})

		Convey("When env var is not set", func() {
			So(envutil.GetDuration("UNSET_DUR_KEY", time.Minute), ShouldEqual, time.Minute)
		})

		Convey("When env var is not a valid duration", func() {
			t.Setenv("BAD_DUR", "not-a-duration")
			So(envutil.GetDuration("BAD_DUR", 5*time.Second), ShouldEqual, 5*time.Second)
		})
	})
}

func TestRequire(t *testing.T) {
	Convey("Given Require", t, func() {
		Convey("When all keys are present", func() {
			t.Setenv("REQ_A", "a")
			t.Setenv("REQ_B", "b")
			err := envutil.Require("REQ_A", "REQ_B")
			So(err, ShouldBeNil)
		})

		Convey("When one key is missing", func() {
			t.Setenv("REQ_C", "c")
			err := envutil.Require("REQ_C", "MISSING_KEY_XYZ")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "MISSING_KEY_XYZ")
		})

		Convey("When multiple keys are missing", func() {
			err := envutil.Require("MISSING_X", "MISSING_Y")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "MISSING_X")
			So(err.Error(), ShouldContainSubstring, "MISSING_Y")
		})
	})
}
