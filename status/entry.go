package status

import (
	"time"

	"github.com/mbolis/mogo/phase"
	"github.com/mbolis/mogo/sign"
)

type Entry struct {
	Date  time.Time
	Time  time.Time
	Phase phase.Phase
	Sign  sign.Sign
}

type Status int

const (
	VeryNegative Status = -2
	Negative     Status = -1
	Neutral      Status = 0
	Positive     Status = +1
	VeryPositive Status = +2
	Warning      Status = 11
)

func (e Entry) Haircut() Status {
	switch e.Sign {
	case sign.Leo, sign.Virgo:
		return Positive

	case sign.Capricorn:
		if e.Phase.IsWaning() {
			return Negative
		}

	case sign.Cancer, sign.Pisces:
		return VeryNegative
	}
	return Neutral
}

func (e Entry) NailsCut() (status Status) {
	switch e.Sign {
	case sign.Cancer, sign.Gemini, sign.Pisces:
		status--
	}

	switch e.Date.Weekday() {
	case time.Friday:
		status++
		if status == 0 {
			return Warning
		}
	case time.Saturday:
		status--
	}

	return
}

func (e Entry) Epilation() (status Status) {
	switch {
	case e.Phase.IsWaning():
		status++
	case e.Phase.IsWaxing():
		status--
	}

	switch e.Sign {
	case sign.Capricorn:
		if e.Phase.IsWaning() {
			status++
		}
	case sign.Leo, sign.Virgo:
		status--
		if status == 0 {
			return Warning
		}
	}

	return
}

func (e Entry) FacialCleansing() (status Status) {
	switch e.Phase {
	case phase.Waxing1, phase.Waxing2, phase.Waxing3:
		status--
		if e.Sign == sign.Leo {
			status--
		}

	case phase.Full:
		status = VeryNegative

	case phase.Waning1, phase.Waning2, phase.Waning3:
		switch e.Sign {
		case sign.Aries, sign.Capricorn:
			status++
		}
	}
	return
}

func (e Entry) FaceMask() (status Status) {
	if e.Phase.IsWaxing() {
		status++
	}
	if e.Sign == sign.Aries {
		status++
	}
	return
}
