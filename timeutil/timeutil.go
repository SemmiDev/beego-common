// Package timeutil provides time formatting, parsing, and comparison helpers
// with special support for Indonesian (WIB/WITA/WIT) time zones and locale-
// aware formatting.
//
// Usage:
//
//	now := time.Now()
//	timeutil.FormatIndonesian(now)     // "5 Maret 2026, 14:01 WIB"
//	timeutil.FormatDate(now)           // "2026-03-05"
//	timeutil.FormatDateTime(now)       // "2026-03-05 14:01:05"
//	timeutil.StartOfDay(now)           // 2026-03-05 00:00:00
//	timeutil.IsExpired(exp)            // true/false
package timeutil

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Indonesian time zones
// ---------------------------------------------------------------------------

var (
	// WIB is Western Indonesian Time (UTC+7) — Jakarta, Bandung, Surabaya.
	WIB = time.FixedZone("WIB", 7*60*60)

	// WITA is Central Indonesian Time (UTC+8) — Bali, Makassar.
	WITA = time.FixedZone("WITA", 8*60*60)

	// WIT is Eastern Indonesian Time (UTC+9) — Papua, Maluku.
	WIT = time.FixedZone("WIT", 9*60*60)
)

// indonesianMonths maps Go month numbers to Indonesian month names.
var indonesianMonths = map[time.Month]string{
	time.January:   "Januari",
	time.February:  "Februari",
	time.March:     "Maret",
	time.April:     "April",
	time.May:       "Mei",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "Agustus",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Desember",
}

// indonesianDays maps Go weekday numbers to Indonesian day names.
var indonesianDays = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

// ---------------------------------------------------------------------------
// Indonesian formatting
// ---------------------------------------------------------------------------

// FormatIndonesian formats a time as "5 Maret 2026, 14:01 WIB".
// The time is converted to WIB before formatting.
func FormatIndonesian(t time.Time) string {
	t = t.In(WIB)
	return fmt.Sprintf("%d %s %d, %02d:%02d WIB",
		t.Day(), indonesianMonths[t.Month()], t.Year(),
		t.Hour(), t.Minute())
}

// FormatIndonesianFull formats with day name: "Senin, 5 Maret 2026, 14:01 WIB".
func FormatIndonesianFull(t time.Time) string {
	t = t.In(WIB)
	return fmt.Sprintf("%s, %d %s %d, %02d:%02d WIB",
		indonesianDays[t.Weekday()],
		t.Day(), indonesianMonths[t.Month()], t.Year(),
		t.Hour(), t.Minute())
}

// FormatIndonesianDate formats a time as "5 Maret 2026" (date only).
func FormatIndonesianDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d",
		t.Day(), indonesianMonths[t.Month()], t.Year())
}

// MonthName returns the Indonesian name for a Go time.Month.
func MonthName(m time.Month) string {
	return indonesianMonths[m]
}

// DayName returns the Indonesian name for a Go time.Weekday.
func DayName(d time.Weekday) string {
	return indonesianDays[d]
}

// ---------------------------------------------------------------------------
// Standard formatting
// ---------------------------------------------------------------------------

// FormatDate formats a time as "2006-01-02" (ISO date).
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime formats a time as "2006-01-02 15:04:05".
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatRFC3339 formats a time as RFC 3339 (e.g. "2006-01-02T15:04:05Z07:00").
func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ---------------------------------------------------------------------------
// Parsing
// ---------------------------------------------------------------------------

// ParseDate parses a "2006-01-02" date string.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// ParseDateTime parses a "2006-01-02 15:04:05" datetime string.
func ParseDateTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// MustParseDate parses a "2006-01-02" date string, panicking on error.
// Use only for constants or test setup.
func MustParseDate(s string) time.Time {
	t, err := ParseDate(s)
	if err != nil {
		panic(fmt.Sprintf("timeutil: invalid date %q: %v", s, err))
	}
	return t
}

// ---------------------------------------------------------------------------
// Boundaries
// ---------------------------------------------------------------------------

// StartOfDay returns t with hour, minute, second, and nanosecond set to zero.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns t with time set to 23:59:59.999999999.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-1), t.Location())
}

// StartOfMonth returns the first day of t's month at 00:00:00.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last day of t's month at 23:59:59.
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, -1)
}

// ---------------------------------------------------------------------------
// Comparison
// ---------------------------------------------------------------------------

// IsExpired reports whether t is in the past.
func IsExpired(t time.Time) bool {
	return time.Now().After(t)
}

// IsFuture reports whether t is in the future.
func IsFuture(t time.Time) bool {
	return t.After(time.Now())
}

// IsSameDay reports whether a and b fall on the same calendar day.
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// DaysBetween returns the number of calendar days between a and b.
// Always returns a non-negative value.
func DaysBetween(a, b time.Time) int {
	a = StartOfDay(a)
	b = StartOfDay(b)
	diff := b.Sub(a)
	days := int(diff.Hours() / 24)
	if days < 0 {
		return -days
	}
	return days
}

// ---------------------------------------------------------------------------
// Convenience
// ---------------------------------------------------------------------------

// NowWIB returns the current time in WIB timezone.
func NowWIB() time.Time {
	return time.Now().In(WIB)
}

// NowPtr returns a pointer to the current time (useful for "updated_at" fields).
func NowPtr() *time.Time {
	now := time.Now()
	return &now
}

// ---------------------------------------------------------------------------
// Calendar
// ---------------------------------------------------------------------------

// StartOfYear returns the first instant of t's year (Jan 1 00:00:00) in t's location.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the last instant of t's year (Dec 31 23:59:59.999999999) in t's location.
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.December, 31, 23, 59, 59, int(time.Second-1), t.Location())
}

// Quarter returns the calendar quarter (1–4) for t.
//
//	timeutil.Quarter(time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)) → 2
func Quarter(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

// WeekRange returns (monday, sunday) for the ISO week containing t.
// Both returned times are at 00:00:00 in t's location.
func WeekRange(t time.Time) (monday, sunday time.Time) {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday → 7
	}
	monday = StartOfDay(t.AddDate(0, 0, -(weekday - 1)))
	sunday = StartOfDay(t.AddDate(0, 0, 7-weekday))
	return
}

// ---------------------------------------------------------------------------
// Age & Duration
// ---------------------------------------------------------------------------

// Age returns the number of complete years from birthDate to today (UTC).
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now().UTC())
}

// AgeAt returns the number of complete years from birthDate to ref.
func AgeAt(birthDate, ref time.Time) int {
	years := ref.Year() - birthDate.Year()
	// Subtract 1 if the birthday hasn't occurred yet this year.
	if ref.Month() < birthDate.Month() ||
		(ref.Month() == birthDate.Month() && ref.Day() < birthDate.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

// FormatDurationIndonesian formats a duration as a human-readable Indonesian string.
//
//	timeutil.FormatDurationIndonesian(2*time.Hour + 15*time.Minute) → "2 jam 15 menit"
//	timeutil.FormatDurationIndonesian(3 * 24 * time.Hour)           → "3 hari"
func FormatDurationIndonesian(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	switch {
	case days > 0 && hours > 0:
		return fmt.Sprintf("%d hari %d jam", days, hours)
	case days > 0:
		return fmt.Sprintf("%d hari", days)
	case hours > 0 && minutes > 0:
		return fmt.Sprintf("%d jam %d menit", hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%d jam", hours)
	case minutes > 0 && seconds > 0:
		return fmt.Sprintf("%d menit %d detik", minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%d menit", minutes)
	default:
		return fmt.Sprintf("%d detik", seconds)
	}
}

// ---------------------------------------------------------------------------
// Convenience (additional)
// ---------------------------------------------------------------------------

// ToPtr returns a pointer to t. Complement to the existing NowPtr.
//
//	ptr := timeutil.ToPtr(someTime)
func ToPtr(t time.Time) *time.Time {
	return &t
}

// ParseDateTimeMulti tries each layout in order and returns the first
// successful parse. Returns an error if none of the layouts match.
//
//	t, err := timeutil.ParseDateTimeMulti(s,
//	    "2006-01-02",
//	    "2006-01-02 15:04:05",
//	    time.RFC3339,
//	)
func ParseDateTimeMulti(s string, layouts ...string) (time.Time, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timeutil: %q does not match any of the provided layouts", s)
}
