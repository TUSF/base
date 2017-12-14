//This package displays the time using dozenal digits.
//
//The format is Hour:Minute:Second:Third, on a radix of 20:50:50;50 (in or 24:60:60:60 in decimal).
package main

import (
	"fmt"
	"time"

	"github.com/TUSF/dozenal"
)

const (
	//	A 60th of a second.
	Third time.Duration = time.Second / 60
)

func main() {
	c := time.Tick(Third)
	for now := range c {
		hour, minute, second := now.Clock()
		third := now.Nanosecond() / int(Third)
		fmt.Printf("\r%02s:%02s:%02s;%02s",
			dozenal.Amer.Int64(int64(hour)),
			dozenal.Amer.Int64(int64(minute)),
			dozenal.Amer.Int64(int64(second)),
			dozenal.Amer.Int64(int64(third)))
	}
}
