package position

import (
	"time"

	"github.com/mbolis/mogo/jd"
	"github.com/mbolis/mogo/util"
	"github.com/mshafiee/swephgo"
)

type Position struct {
	Longitude float64
	Latitude  float64
	Distance  float64
	JD        float64
}

func Calc(jd float64, planet int) Position {
	return calc(jd, false, planet)
}

func CalcUT(jd float64, planet int) Position {
	return calc(jd, true, planet)
}

func CalcTime(d time.Time, planet int) Position {
	jd := jd.FromTime(d)
	return calc(jd, false, planet)
}

func calc(jd float64, ut bool, planet int) Position {
	var calc func(jd float64, pl int, flag int, xx []float64, err []byte) int32
	if ut {
		calc = swephgo.CalcUt
	} else {
		calc = swephgo.Calc
	}

	var xx [6]float64
	var errMsg [256]byte
	result := calc(jd, planet, 0, xx[:], errMsg[:])
	if result == swephgo.Err {
		panic(util.NTString(errMsg[:]))
	}

	return Position{
		Longitude: xx[0],
		Latitude:  xx[1],
		Distance:  xx[2],
		JD:        jd,
	}
}
