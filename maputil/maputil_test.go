package maputil_test

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/beego-common/maputil"
)

func TestKeys(t *testing.T) {
	Convey("Given Keys", t, func() {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		keys := maputil.Keys(m)
		sort.Strings(keys)
		So(keys, ShouldResemble, []string{"a", "b", "c"})
	})
}

func TestValues(t *testing.T) {
	Convey("Given Values", t, func() {
		m := map[string]int{"a": 1}
		vals := maputil.Values(m)
		So(len(vals), ShouldEqual, 1)
		So(vals[0], ShouldEqual, 1)
	})
}

func TestMerge(t *testing.T) {
	Convey("Given Merge", t, func() {
		Convey("When later map wins on conflict", func() {
			a := map[string]int{"x": 1, "y": 2}
			b := map[string]int{"y": 99, "z": 3}
			result := maputil.Merge(a, b)
			So(result["x"], ShouldEqual, 1)
			So(result["y"], ShouldEqual, 99)
			So(result["z"], ShouldEqual, 3)
		})

		Convey("When zero args returns empty non-nil map", func() {
			result := maputil.Merge[string, int]()
			So(result, ShouldNotBeNil)
			So(len(result), ShouldEqual, 0)
		})
	})
}

func TestFilter(t *testing.T) {
	Convey("Given Filter", t, func() {
		m := map[string]int{"a": 1, "b": -1, "c": 2}
		positive := maputil.Filter(m, func(_ string, v int) bool { return v > 0 })
		So(len(positive), ShouldEqual, 2)
		So(positive["a"], ShouldEqual, 1)
		So(positive["c"], ShouldEqual, 2)
	})
}

func TestMapValues(t *testing.T) {
	Convey("Given MapValues", t, func() {
		m := map[string]string{"hello": "world"}
		lengths := maputil.MapValues(m, func(v string) int { return len(v) })
		So(lengths["hello"], ShouldEqual, 5)
	})
}

func TestInvert(t *testing.T) {
	Convey("Given Invert", t, func() {
		m := map[string]int{"a": 1, "b": 2}
		inv := maputil.Invert(m)
		So(inv[1], ShouldEqual, "a")
		So(inv[2], ShouldEqual, "b")
	})
}

func TestPick(t *testing.T) {
	Convey("Given Pick", t, func() {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		picked := maputil.Pick(m, "a", "c")
		So(len(picked), ShouldEqual, 2)
		So(picked["a"], ShouldEqual, 1)
		So(picked["c"], ShouldEqual, 3)
		_, exists := picked["b"]
		So(exists, ShouldBeFalse)
	})
}

func TestOmit(t *testing.T) {
	Convey("Given Omit", t, func() {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		omitted := maputil.Omit(m, "b")
		So(len(omitted), ShouldEqual, 2)
		_, exists := omitted["b"]
		So(exists, ShouldBeFalse)
	})
}

func TestContains(t *testing.T) {
	Convey("Given Contains", t, func() {
		m := map[string]int{"a": 1}
		So(maputil.Contains(m, "a"), ShouldBeTrue)
		So(maputil.Contains(m, "z"), ShouldBeFalse)
	})
}

func TestGetOrDefault(t *testing.T) {
	Convey("Given GetOrDefault", t, func() {
		m := map[string]int{"a": 1}

		Convey("When key exists", func() {
			So(maputil.GetOrDefault(m, "a", 99), ShouldEqual, 1)
		})

		Convey("When key is missing", func() {
			So(maputil.GetOrDefault(m, "missing", 99), ShouldEqual, 99)
		})

		Convey("When map is nil", func() {
			So(maputil.GetOrDefault[string, int](nil, "k", 42), ShouldEqual, 42)
		})
	})
}
