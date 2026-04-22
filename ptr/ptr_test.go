package ptr_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/beego-common/ptr"
)

// ---------------------------------------------------------------------------
// Value → Pointer
// ---------------------------------------------------------------------------

func TestValueToPointer(t *testing.T) {
	Convey("Given the value-to-pointer helpers", t, func() {
		Convey("String should return a pointer to the string", func() {
			p := ptr.String("hello")
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, "hello")
		})

		Convey("Int should return a pointer to the int", func() {
			p := ptr.Int(42)
			So(*p, ShouldEqual, 42)
		})

		Convey("Int32 should return a pointer to the int32", func() {
			p := ptr.Int32(32)
			So(*p, ShouldEqual, int32(32))
		})

		Convey("Int64 should return a pointer to the int64", func() {
			p := ptr.Int64(64)
			So(*p, ShouldEqual, int64(64))
		})

		Convey("Float32 should return a pointer to the float32", func() {
			p := ptr.Float32(3.14)
			So(*p, ShouldAlmostEqual, float32(3.14))
		})

		Convey("Float64 should return a pointer to the float64", func() {
			p := ptr.Float64(2.718)
			So(*p, ShouldAlmostEqual, 2.718)
		})

		Convey("Bool should return a pointer to the bool", func() {
			p := ptr.Bool(true)
			So(*p, ShouldBeTrue)
		})

		Convey("Time should return a pointer to time.Time", func() {
			now := time.Now()
			p := ptr.Time(now)
			So(*p, ShouldEqual, now)
		})

		Convey("Of (generic) should work with any type", func() {
			sp := ptr.Of("generic")
			So(*sp, ShouldEqual, "generic")

			ip := ptr.Of(99)
			So(*ip, ShouldEqual, 99)
		})
	})
}

// ---------------------------------------------------------------------------
// Pointer → Value (safe dereference)
// ---------------------------------------------------------------------------

func TestPointerToValue(t *testing.T) {
	Convey("Given the pointer-to-value helpers", t, func() {
		Convey("StringValue should return the value or empty string", func() {
			s := "hello"
			So(ptr.StringValue(&s), ShouldEqual, "hello")
			So(ptr.StringValue(nil), ShouldBeEmpty)
		})

		Convey("IntValue should return the value or 0", func() {
			i := 42
			So(ptr.IntValue(&i), ShouldEqual, 42)
			So(ptr.IntValue(nil), ShouldEqual, 0)
		})

		Convey("Int64Value should return the value or 0", func() {
			i := int64(64)
			So(ptr.Int64Value(&i), ShouldEqual, int64(64))
			So(ptr.Int64Value(nil), ShouldEqual, int64(0))
		})

		Convey("Float64Value should return the value or 0", func() {
			f := 3.14
			So(ptr.Float64Value(&f), ShouldAlmostEqual, 3.14)
			So(ptr.Float64Value(nil), ShouldEqual, 0.0)
		})

		Convey("BoolValue should return the value or false", func() {
			b := true
			So(ptr.BoolValue(&b), ShouldBeTrue)
			So(ptr.BoolValue(nil), ShouldBeFalse)
		})

		Convey("TimeValue should return the value or zero time", func() {
			now := time.Now()
			So(ptr.TimeValue(&now), ShouldEqual, now)
			So(ptr.TimeValue(nil), ShouldResemble, time.Time{})
		})

		Convey("ValueOf (generic) should return the value or zero", func() {
			s := "hi"
			So(ptr.ValueOf(&s), ShouldEqual, "hi")
			So(ptr.ValueOf[string](nil), ShouldBeEmpty)
		})

		Convey("ValueOrDefault should return value or the provided default", func() {
			s := "present"
			So(ptr.ValueOrDefault(&s, "fallback"), ShouldEqual, "present")
			So(ptr.ValueOrDefault[string](nil, "fallback"), ShouldEqual, "fallback")
		})
	})
}

// ---------------------------------------------------------------------------
// Nil checks
// ---------------------------------------------------------------------------

func TestIsNilOrEmpty(t *testing.T) {
	Convey("Given IsNilOrEmpty", t, func() {
		Convey("When pointer is nil", func() {
			So(ptr.IsNilOrEmpty(nil), ShouldBeTrue)
		})

		Convey("When pointer points to empty string", func() {
			s := ""
			So(ptr.IsNilOrEmpty(&s), ShouldBeTrue)
		})

		Convey("When pointer points to non-empty string", func() {
			s := "hello"
			So(ptr.IsNilOrEmpty(&s), ShouldBeFalse)
		})
	})
}

// ---------------------------------------------------------------------------
// Value → Pointer (additional numeric types)
// ---------------------------------------------------------------------------

func TestValueToPointerAdditional(t *testing.T) {
	Convey("Given additional value-to-pointer helpers", t, func() {
		Convey("Uint should return a non-nil pointer with correct value", func() {
			p := ptr.Uint(uint(10))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, uint(10))
		})

		Convey("Uint32 should return a non-nil pointer with correct value", func() {
			p := ptr.Uint32(uint32(32))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, uint32(32))
		})

		Convey("Uint64 should return a non-nil pointer with correct value", func() {
			p := ptr.Uint64(uint64(64))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, uint64(64))
		})

		Convey("Int8 should return a non-nil pointer with correct value", func() {
			p := ptr.Int8(int8(8))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, int8(8))
		})

		Convey("Int16 should return a non-nil pointer with correct value", func() {
			p := ptr.Int16(int16(16))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, int16(16))
		})

		Convey("Byte should return a non-nil pointer with correct value", func() {
			p := ptr.Byte(byte('A'))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, byte('A'))
		})

		Convey("Rune should return a non-nil pointer with correct value", func() {
			p := ptr.Rune(rune('Z'))
			So(p, ShouldNotBeNil)
			So(*p, ShouldEqual, rune('Z'))
		})
	})
}

// ---------------------------------------------------------------------------
// Pointer → Value (additional types)
// ---------------------------------------------------------------------------

func TestPointerToValueAdditional(t *testing.T) {
	Convey("Given additional pointer-to-value helpers", t, func() {
		Convey("Int32Value should return the value or 0", func() {
			v := int32(32)
			So(ptr.Int32Value(&v), ShouldEqual, int32(32))
			So(ptr.Int32Value(nil), ShouldEqual, int32(0))
		})

		Convey("Float32Value should return the value or 0", func() {
			v := float32(3.14)
			So(ptr.Float32Value(&v), ShouldAlmostEqual, float32(3.14))
			So(ptr.Float32Value(nil), ShouldEqual, float32(0))
		})

		Convey("UintValue should return the value or 0", func() {
			v := uint(10)
			So(ptr.UintValue(&v), ShouldEqual, uint(10))
			So(ptr.UintValue(nil), ShouldEqual, uint(0))
		})

		Convey("Uint32Value should return the value or 0", func() {
			v := uint32(32)
			So(ptr.Uint32Value(&v), ShouldEqual, uint32(32))
			So(ptr.Uint32Value(nil), ShouldEqual, uint32(0))
		})

		Convey("Uint64Value should return the value or 0", func() {
			v := uint64(64)
			So(ptr.Uint64Value(&v), ShouldEqual, uint64(64))
			So(ptr.Uint64Value(nil), ShouldEqual, uint64(0))
		})

		Convey("Int8Value should return the value or 0", func() {
			v := int8(8)
			So(ptr.Int8Value(&v), ShouldEqual, int8(8))
			So(ptr.Int8Value(nil), ShouldEqual, int8(0))
		})

		Convey("Int16Value should return the value or 0", func() {
			v := int16(16)
			So(ptr.Int16Value(&v), ShouldEqual, int16(16))
			So(ptr.Int16Value(nil), ShouldEqual, int16(0))
		})
	})
}

