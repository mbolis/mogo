package icons

import (
	"fmt"
	"strings"

	"github.com/mbolis/mogo/phase"
	"github.com/mbolis/mogo/sign"
	"github.com/mbolis/mogo/status"
)

type Style int

const (
	Arrows Style = iota
	Thumbs
	Semaphore
)

func (style Style) positive() rune {
	switch style {
	case Arrows:
		return '🔼'
	case Thumbs:
		return '👍'
	case Semaphore:
		return '🟢'
	default:
		panic("unrecognized style")
	}
}

func (style Style) negative() rune {
	switch style {
	case Arrows:
		return '🔽'
	case Thumbs:
		return '👎'
	case Semaphore:
		return '🔴'
	default:
		panic("unrecognized style")
	}
}

func (style Style) warning() rune {
	switch style {
	case Arrows:
		return '🔁'
	case Thumbs:
		return '✋'
	case Semaphore:
		return '🟡'
	default:
		panic("unrecognized style")
	}
}

func (style Style) Status(s status.Status) string {
	switch {
	case s == status.Warning:
		return string(style.warning())
	case s < 0:
		return strings.Repeat(string(style.negative()), -int(s))
	case s > 0:
		return strings.Repeat(string(style.positive()), int(s))
	default:
		return ""
	}
}

func (Style) Sign(s sign.Sign) rune {
	switch s {
	case sign.Aries:
		return '♈'
	case sign.Taurus:
		return '♉'
	case sign.Gemini:
		return '♊'
	case sign.Cancer:
		return '♋'
	case sign.Leo:
		return '♌'
	case sign.Virgo:
		return '♍'
	case sign.Libra:
		return '♎'
	case sign.Scorpio:
		return '♏'
	case sign.Sagittarius:
		return '♐'
	case sign.Capricorn:
		return '♑'
	case sign.Aquarius:
		return '♒'
	case sign.Pisces:
		return '♓'
	default:
		panic(fmt.Sprintf("unknown sign: %d", s))
	}
}

func (Style) Phase(ph phase.Phase) rune {
	switch ph {
	case phase.New:
		return '🌑'
	case phase.Waxing1:
		return '🌒'
	case phase.Waxing2:
		return '🌓'
	case phase.Waxing3:
		return '🌔'
	case phase.Full:
		return '🌕'
	case phase.Waning1:
		return '🌖'
	case phase.Waning2:
		return '🌗'
	case phase.Waning3:
		return '🌘'
	}
	panic("impossible moon phase")
}
