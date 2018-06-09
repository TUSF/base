// A demonstration of the TUSF/base package, meant for the web browser.
//
// Got to https://TUSF.github.io/base
package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/TUSF/base"
	"github.com/TUSF/base/dozenal"
	"honnef.co/go/js/dom"
)

//Currently in-use Formatter.
var BAS base.Formatter

//Optional formats.
//Base 6
var sex base.Formatter = base.NewFormatter([]string{"0", "1", "2", "3", "4", "5"})

//Base 36 (6^2)
var hexsex base.Formatter = base.NewFormatter([]string{"0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "A", "B",
	"C", "D", "E", "F", "G", "H",
	"I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z"})

//Base 16 (hexadecimal)
var hex base.Formatter = base.NewFormatter([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F"})

//Base 20 Vigesimal, based on Maya Glyphs
var maya base.Formatter = base.NewFormatter([]string{
	"\u1d2e0", "\u1d2e1", "\u1d2e2", "\u1d2e3", "\u1d2e4",
	"\u1d2e5", "\u1d2e6", "\u1d2e7", "\u1d2e8", "\u1d2e9",
	"\u1d2ea", "\u1d2eb", "\u1d2ec", "\u1d2ed", "\u1d2ee",
	"\u1d2ef", "\u1d2f0", "\u1d2f1", "\u1d2f2", "\u1d2f3",
})

func main() {
	document := dom.GetWindow().Document()
	input := document.GetElementByID("input").(*dom.HTMLInputElement)
	output := document.GetElementByID("output")
	options := document.GetElementByID("base").(*dom.HTMLSelectElement)
	BAS = dozenal.ASCII

	input.AddEventListener("keydown", false, func(e dom.Event) {
		ke := e.(*dom.KeyboardEvent)
		if ke.KeyCode == '\r' {
			if input.Value != "" {
				output.SetTextContent(Convert(input.Value))
			}
			ke.PreventDefault()
		}
	})
	options.AddEventListener("change", false, func(e dom.Event) {
		switch options.Value {
		default:
			BAS = dozenal.ASCII
		case "dozamer":
			BAS = dozenal.Amer
		case "dozbrit":
			BAS = dozenal.Brit
		case "sex":
			BAS = sex
		case "hexsex":
			BAS = hexsex
		case "hex":
			BAS = hex
		case "maya":
			BAS = maya
		}
		if input.Value != "" {
			output.SetTextContent(Convert(input.Value))
		}
	})
}

// Take the string, and convert it into a string of the active base.
func Convert(s string) string {
	var INT big.Int
	var RAT big.Rat

	s = strings.TrimSpace(s)

	if INT, t := INT.SetString(s, 0); t {
		//First, assume it's a plain integer.
		return BAS.BigInt(INT)

	} else if RAT, t := RAT.SetString(s); t {
		//Second, assume it's a fraction. ("12/7")
		return BAS.BigRat(RAT)

	} else {
		//Third, assume it's a decimal number. ("10.123")
		if strings.Index(s, ".") > -1 {
			if nums := strings.Split(s, "."); len(nums) == 2 {
				//Convert each side of the decimal point into a big.Int?
				if nums[1] == "" {
					return "Not a valid number. Integers, Fractions or Decimals only."
				}
				if nums[0] == "" {
					RAT.SetInt64(0)
				} else {
					if _, err := strconv.Atoi(nums[0]); err != nil {
						return "Not a valid number. Integers, Fractions or Decimals only."
					}
					if _, err := strconv.Atoi(nums[1]); err != nil {
						return "Not a valid number. Integers, Fractions or Decimals only."
					}
					RAT, t := RAT.SetString(nums[0])
					if !t {
						return "Not a valid number. Integers, Fractions or Decimals only."
					}

					d, t := new(big.Rat).SetString(
						fmt.Sprintf("%s/1%s", nums[1], strings.Repeat("0", len(nums))))
					if !t {
						return "Not a valid number. Integers, Fractions or Decimals only."
					}
					RAT.Add(RAT, d)
					return BAS.BigRat(RAT)
				}
			} else {
				return "Not a valid number. Integers, Fractions or Decimals only."
			}
		} else {
			return "Not a valid number. Integers, Fractions or Decimals only."
		}
	}
	return ""
}
