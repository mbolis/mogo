package model

import (
	"fmt"
	"time"
)

type DailyValue[T ~int] struct {
	Curr  T
	Next  T
	Event *Event[T]
}

func (dv DailyValue[T]) String() string {
	out := fmt.Sprintf("%v", dv.Curr)
	if dv.Curr != dv.Next {
		out += fmt.Sprintf("->%s->%v", dv.Event, dv.Next)
	}
	return out
}

func (dv DailyValue[T]) Value() T {
	if dv.Event != nil {
		return dv.Event.Value
	}
	return dv.Curr
}

type Event[T ~int] struct {
	Time  time.Time
	Value T
}

func (e *Event[T]) String() string {
	return fmt.Sprintf("%v (%s)", e.Value, e.Time.Format("15:04"))
}
