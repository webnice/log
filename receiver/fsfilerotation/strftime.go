package fsfilerotation // import "github.com/webnice/log/v2/receiver/fsfilerotation"

/*
Example:
	str, err := strftime.Format("%d.%m.%Y", time.Now()) // 30.12.2016

Directives:
	%a - Locale’s abbreviated weekday name
	%A - Locale’s full weekday name
	%b - Locale’s abbreviated month name
	%B - Locale’s full month name
	%c - Locale’s appropriate date and time representation
	%d - Day of the month as a decimal number [01,31]
	%H - Hour (24-hour clock) as a decimal number [00,23]
	%I - Hour (12-hour clock) as a decimal number [01,12]
	%j - Day of year
	%m - Month as a decimal number [01,12]
	%M - Minute as a decimal number [00,59]
	%p - Locale’s equivalent of either AM or PM
	%S - Second as a decimal number [00,61]
	%U - Week number of the year
	%w - Weekday as a decimal number
	%W - Week number of the year
	%x - Locale’s appropriate date representation
	%X - Locale’s appropriate time representation
	%y - Year without century as a decimal number [00,99]
	%Y - Year with century as a decimal number
	%Z - Time zone name (no characters if no time zone exists)
*/

import (
	"fmt"
	"regexp"
	"time"
)

const (
	timeWeek = time.Hour * 24 * 7
)

var (
	// Map for golang template
	conv = map[string]string{
		"%a": "Mon",        // Locale’s abbreviated weekday name
		"%A": "Monday",     // Locale’s full weekday name
		"%b": "Jan",        // Locale’s abbreviated month name
		"%B": "January",    // Locale’s full month name
		"%c": time.RFC1123, // Locale’s appropriate date and time representation
		"%d": "02",         // Day of the month as a decimal number [01,31]
		"%H": "15",         // Hour (24-hour clock) as a decimal number [00,23]
		"%I": "3",          // Hour (12-hour clock) as a decimal number [01,12]
		"%m": "01",         // Month as a decimal number [01,12]
		"%M": "04",         // Minute as a decimal number [00,59]
		"%p": "PM",         // Locale’s equivalent of either AM or PM
		"%S": "05",         // Second as a decimal number [00,61]
		"%x": "01/02/06",   // Locale’s appropriate date representation
		"%X": "15:04:05",   // Locale’s appropriate time representation
		"%y": "06",         // Year without century as a decimal number [00,99]
		"%Y": "2006",       // Year with century as a decimal number
		"%Z": "MST",        // Time zone name (no characters if no time zone exists)
	}
	rexFmt = regexp.MustCompile("%[%a-zA-Z]")
)

// replaces % directives with right time, will panic on unknown directive
func replaces(match string, t time.Time) (ret string, err error) {
	var (
		format    string
		ok        bool
		st        time.Time
		day, week int
	)

	if match == "%%" {
		ret = "%"
		return
	}
	if format, ok = conv[match]; ok {
		ret = t.Format(format)
		return
	}
	switch match {
	case "%j":
		st = time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
		day = int(t.Sub(st).Hours()/24) + 1
		ret = fmt.Sprintf("%03d", day)
		return
	case "%w":
		ret = fmt.Sprintf("%d", t.Weekday())
		return
	case "%W", "%U":
		st = time.Date(t.Year(), time.January, 1, 23, 0, 0, 0, time.UTC)
		week = 0
		for st.Before(t) {
			week++
			st = st.Add(timeWeek)
		}
		ret = fmt.Sprintf("%02d", week)
		return
	}
	err = fmt.Errorf("Unknown directive: %s", match)

	return
}

// Format return string with % directives expanded.
// Will return error on unknown directive.
func Format(format string, t time.Time) (ret string, err error) {
	defer func() {
		if e := recover(); e != nil {
			// error already set
		}
	}()
	ret = rexFmt.ReplaceAllStringFunc(format, func(match string) (ret string) {
		if ret, err = replaces(match, t); err != nil {
			panic(`break`)
		}
		return
	})

	return
}
