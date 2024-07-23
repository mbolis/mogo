package jd

import (
	"math"
	"time"

	"github.com/mbolis/mogo/util"
	"github.com/mshafiee/swephgo"
	"github.com/soniakeys/meeus/v3/julian"
)

func FromTime(d time.Time) float64 {
	d = d.In(time.UTC)
	et, _ := fromTime(d)
	return et
}

func FromTimeUT(d time.Time) float64 {
	d = d.In(time.UTC)
	_, ut := fromTime(d)
	return ut
}

func fromTime(d time.Time) (et float64, ut float64) {
	d = d.In(time.UTC)
	var ret [2]float64
	var errMsg [256]byte
	if r := swephgo.UtcToJd(
		d.Year(), int(d.Month()), d.Day(),
		d.Hour(), d.Minute(), float64(d.Second())+float64(d.Nanosecond())/1000000000,
		swephgo.SeGregCal,
		ret[:], errMsg[:],
	); r == swephgo.Err {
		panic(util.NTString(errMsg[:]))
	}
	return ret[0], ret[1]
}

func Time(jd float64) time.Time {
	return toTime(jd, false)
}

func TimeUT(jd float64) time.Time {
	return toTime(jd, true)
}

func toTime(jd float64, ut bool) time.Time {
	var jdToUTC func(jd float64, gregflag int, year []int, month []int, day []int, hour []int, min []int, sec []float64)
	if ut {
		jdToUTC = swephgo.Jdut1ToUtc
	} else {
		jdToUTC = swephgo.JdetToUtc
	}

	var year, month, day [1]int
	var hour, min [1]int
	var sec [1]float64
	jdToUTC(jd, swephgo.SeGregCal, year[:], month[:], day[:], hour[:], min[:], sec[:])

	s, f := math.Modf(sec[0])
	return time.Date(year[0], time.Month(month[0]), day[0], hour[0], min[0], int(s), int(f*1000000000), time.UTC)
}

var HalfMinute = julian.TimeToJD(julian.JDToTime(0).Add(30 * time.Second))
