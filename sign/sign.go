package sign

import (
	"fmt"
	"time"

	"github.com/mbolis/mogo/jd"
	"github.com/mbolis/mogo/model"
	"github.com/mbolis/mogo/position"
	"github.com/mshafiee/swephgo"
)

type Sign int

const (
	Aries Sign = iota
	Taurus
	Gemini
	Cancer
	Leo
	Virgo
	Libra
	Scorpio
	Sagittarius
	Capricorn
	Aquarius
	Pisces
)

func OfPosition(p position.Position) Sign {
	return Sign(p.Longitude / 30)
}

func OfLongitude(lon float64) Sign {
	return Sign(lon / 30)
}

func (s Sign) Prev() Sign {
	return (s + 11) % 12
}

func (s Sign) String() string {
	switch s {
	case Aries:
		return "Aries"
	case Taurus:
		return "Taurus"
	case Gemini:
		return "Gemini"
	case Cancer:
		return "Cancer"
	case Leo:
		return "Leo"
	case Virgo:
		return "Virgo"
	case Libra:
		return "Libra"
	case Scorpio:
		return "Scorpio"
	case Sagittarius:
		return "Sagittarius"
	case Capricorn:
		return "Capricorn"
	case Aquarius:
		return "Aquarius"
	case Pisces:
		return "Pisces"
	default:
		panic(fmt.Sprintf("unknown sign: %d", s))
	}
}

type pos position.Position

func (p pos) Sign() Sign {
	return Sign(p.Longitude / 30)
}

func (p pos) Time() time.Time {
	return jd.Time(p.JD)
}

var positionCache = make(map[float64]pos)

func calcCached(d float64) pos {
	if pos, ok := positionCache[d]; ok {
		return pos
	}

	pos := pos(position.Calc(d, swephgo.SeMoon))
	positionCache[d] = pos
	return pos
}

func ForDay(d time.Time) (dv model.DailyValue[Sign]) {
	d0 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).In(time.UTC)
	pos0 := calcCached(jd.FromTime(d0))
	pos1 := calcCached(jd.FromTime(d0.AddDate(0, 0, 1)))

	dv.Curr = pos0.Sign()
	dv.Next = pos1.Sign()

	if dv.Curr != dv.Next {
		dv.Event = binarySearch(pos0, pos1)
		dv.Event.Time = dv.Event.Time.In(d.Location())
	}

	return
}

func binarySearch(start, end pos) *model.Event[Sign] {
	for {
		mid := calcCached(start.JD + (end.JD-start.JD)/2)

		if start.Sign() != mid.Sign() {
			end = mid
		} else {
			start = mid
		}
		if end.JD-start.JD < jd.HalfMinute {
			return &model.Event[Sign]{
				Time:  end.Time(),
				Value: end.Sign(),
			}
		}
	}
}
