package ctxval_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/beego-common/ctxval"
)

type ctxKey string

const (
	nameKey ctxKey = "name"
	ageKey  ctxKey = "age"
)

func TestSet(t *testing.T) {
	Convey("Given Set", t, func() {
		ctx := ctxval.Set[string](context.Background(), nameKey, "Budi")

		Convey("Then Get should retrieve the value", func() {
			v, ok := ctxval.Get[string](ctx, nameKey)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, "Budi")
		})
	})
}

func TestGet(t *testing.T) {
	Convey("Given Get", t, func() {
		Convey("When key exists with correct type", func() {
			ctx := ctxval.Set[int](context.Background(), ageKey, 25)
			v, ok := ctxval.Get[int](ctx, ageKey)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 25)
		})

		Convey("When key is absent", func() {
			v, ok := ctxval.Get[string](context.Background(), nameKey)
			So(ok, ShouldBeFalse)
			So(v, ShouldBeEmpty)
		})

		Convey("When type assertion fails", func() {
			// Store an int but retrieve as string
			ctx := ctxval.Set[int](context.Background(), ageKey, 42)
			v, ok := ctxval.Get[string](ctx, ageKey)
			So(ok, ShouldBeFalse)
			So(v, ShouldBeEmpty)
		})
	})
}

func TestMustGet(t *testing.T) {
	Convey("Given MustGet", t, func() {
		Convey("When key exists", func() {
			ctx := ctxval.Set[string](context.Background(), nameKey, "Siti")
			v := ctxval.MustGet[string](ctx, nameKey)
			So(v, ShouldEqual, "Siti")
		})

		Convey("When key is absent it should panic", func() {
			So(func() {
				ctxval.MustGet[string](context.Background(), nameKey)
			}, ShouldPanic)
		})
	})
}

func TestGetOrDefault(t *testing.T) {
	Convey("Given GetOrDefault", t, func() {
		Convey("When key exists", func() {
			ctx := ctxval.Set[string](context.Background(), nameKey, "Ahmad")
			v := ctxval.GetOrDefault[string](ctx, nameKey, "unknown")
			So(v, ShouldEqual, "Ahmad")
		})

		Convey("When key is absent", func() {
			v := ctxval.GetOrDefault[string](context.Background(), nameKey, "anonymous")
			So(v, ShouldEqual, "anonymous")
		})
	})
}
