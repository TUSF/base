//This package displays the time using Decimal digits.
//
//The format is Hour:Minute:Second:Third, on a radix of 20:50:50;50 (in or 24:60:60:60 in decimal).
package main

import (
	"fmt"
	"time"

	"github.com/TUSF/doz"
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
			doz.Str(hour), doz.Str(minute), doz.Str(second), doz.Str(third))
	}
}
