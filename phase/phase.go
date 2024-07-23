package phase

import (
	"fmt"
	"math"
	"time"

	"github.com/mbolis/mogo/jd"
	"github.com/mbolis/mogo/model"
	"github.com/mbolis/mogo/position"
	"github.com/mshafiee/swephgo"
)

type Value struct {
	Ph, JD float64
}

func (e Value) Time() time.Time {
	return jd.Time(e.JD)
}

func (e Value) Phase() Phase {
	switch {
	case e.Ph == 0:
		return New
	case 0 < e.Ph && e.Ph <= 60:
		return Waxing1
	case 60 < e.Ph && e.Ph <= 120:
		return Waxing2
	case 120 < e.Ph && e.Ph < 180:
		return Waxing3
	case e.Ph == 180:
		return Full
	case -180 < e.Ph && e.Ph <= -120:
		return Waning1
	case -120 < e.Ph && e.Ph <= -60:
		return Waning2
	case -60 < e.Ph && e.Ph < 0:
		return Waning3
	}
	panic(fmt.Sprintf("impossible phase: %f", e.Ph))
}

func Calc(jd float64) Value {
	return calc(jd, false)
}

func CalcUT(jd float64) Value {
	return calc(jd, true)
}

func CalcTime(d time.Time) Value {
	jd := jd.FromTime(d)
	return calc(jd, false)
}

func calc(jd float64, ut bool) Value {
	var calcPos func(jd float64, planet int) position.Position
	if ut {
		calcPos = position.CalcUT
	} else {
		calcPos = position.Calc
	}

	sunLon := calcPos(jd, swephgo.SeSun).Longitude
	moonLon := calcPos(jd, swephgo.SeMoon).Longitude
	return Value{
		JD: jd,
		Ph: normDeg180(moonLon - sunLon),
	}
}
func normDeg180(th float64) float64 {
	th = math.Mod(th, 360)
	if th < 0 {
		th += 360
	}
	if th < 1e-13 {
		th = 0
	}
	if th >= 180 {
		th -= 360
	}
	return th
}

type Phase int

const (
	New Phase = iota
	Waxing1
	Waxing2
	Waxing3
	Full
	Waning1
	Waning2
	Waning3
)

func (p Phase) IsWaning() bool {
	return p == Waning1 || p == Waning2 || p == Waning3
}

func (p Phase) IsWaxing() bool {
	return p == Waxing1 || p == Waxing2 || p == Waxing3
}

func (p Phase) String() string {
	switch p {
	case New:
		return "New"
	case Waxing1, Waxing2, Waxing3:
		return "Waxing"
	case Full:
		return "Full"
	case Waning1, Waning2, Waning3:
		return "Waning"
	default:
		panic(fmt.Sprintf("unrecognized EventType: %d", p))
	}
}

var phaseCache = make(map[float64]Value)

func calcCached(d float64) Value {
	if ph, ok := phaseCache[d]; ok {
		return ph
	}

	ph := Calc(d)
	phaseCache[d] = ph
	return ph
}

func ForDay(d time.Time) (dv model.DailyValue[Phase]) {
	d0 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).In(time.UTC)
	ph0 := calcCached(jd.FromTime(d0))
	ph1 := calcCached(jd.FromTime(d0.AddDate(0, 0, 1)))

	dv.Curr = ph0.Phase()
	dv.Next = ph1.Phase()

	switch {
	case ph0.Ph <= 0 && ph1.Ph > 0:
		dv.Event = New.binarySearch(ph0, ph1)
		dv.Event.Time = dv.Event.Time.In(d.Location())
	case ph0.Ph > 0 && ph1.Ph < 0:
		dv.Event = Full.binarySearch(ph0, ph1)
		dv.Event.Time = dv.Event.Time.In(d.Location())
	}

	return
}

func (et Phase) binarySearch(start, end Value) *model.Event[Phase] {
	for {
		mid := calcCached(start.JD + (end.JD-start.JD)/2)
		start, end = et.selectNext(start, mid, end)
		if end.JD-start.JD < jd.HalfMinute {
			return &model.Event[Phase]{
				Time:  end.Time(),
				Value: et,
			}
		}
	}
}
func (et Phase) selectNext(start, mid, end Value) (Value, Value) {
	switch et {
	case New:
		switch {
		case mid.Ph < 0:
			return mid, end
		case mid.Ph == 0:
			return mid, mid
		case mid.Ph > 0:
			return start, mid
		}
	case Full:
		switch {
		case mid.Ph < 0:
			return start, mid
		case mid.Ph == 0:
			return mid, mid
		case mid.Ph > 0:
			return mid, end
		}
	}

	panic(fmt.Sprintf("unrecognized EventType: %d", et))
}
