//	This package implements a set of functions useful for displaying numbers as dozenal,
//	otherwise called duodecimal.
//
//	Built-in integer values can make use of Str(), while Int() and Rat() can be used for
//	Int and Rat structs from the math/big package.
package doz

import (
	"math/big"
)

//	Digits is a slice of string, instead of just a string, because Digits[10] should
//	return a whole character; not a byte. Exported, so that anyone can modify any
//	part of the slice for their needs. Be default, ASCIIDigits is used.
var Digits = ASCIIDigits

//	An ASCII-friendly set of Dozenal digits. Dec and El are represented using the Latin
//	letters "X" and "E", to be compatible with most displays.
var ASCIIDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "X", "E"}

//	The set of Dozenal digits recommended by The Dozenal Society of America. Dec and El are
//	represented using the Greek letter, Chi (Χ) and the Latin "open E". (Ɛ)
var AmerDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u03A7", "\u0190"}

//	The set of Dozenal digits recommended by The Dozenal Society of Great Britain . Dec and El are
//	represented using a turned 2 (↊) and 3 (↋). These are the least supported among fonts.
var BritDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u218A", "\u218B"}

//Likely to be initialized a lot, so may as well reuse the same one.
var zero *big.Int = big.NewInt(0)
var dozen *big.Int = big.NewInt(0xC)

//Formats dozenal integers.
func integer(n uint64, signed bool) string {
	if n == 0 {
		return "0"
	}
	var str string

	//This process is simple. Divide the integer by a dozen (0xC), and write the
	// remainder to the string... Well, that's technically what's going on.
	for u := n; u != 0; u /= 0xC {
		rem := u % 0xC
		str = Digits[rem] + str
	}

	if signed && int64(n) < 0 {
		//n = -n
		str = "-" + str
	}
	return str
}

// Converts a built-in integer value into a dozenal string.
// Non-integer variables return a blank string.
func Str(n interface{}) string {
	switch v := n.(type) {
	case int:
		return integer(uint64(v), true)
	case int8:
		return integer(uint64(v), true)
	case int16:
		return integer(uint64(v), true)
	case int32:
		return integer(uint64(v), true)
	case int64:
		return integer(uint64(v), true)
	case uint:
		return integer(uint64(v), false)
	case uint8:
		return integer(uint64(v), false)
	case uint16:
		return integer(uint64(v), false)
	case uint32:
		return integer(uint64(v), false)
	case uint64:
		return integer(v, false)
	default:
		return ""
	}
}

//Accepts a big.Int from the `math/big` package, and returns a dozenal integer as a string.
func Int(i *big.Int) string {
	//No need to waste time if it's 0. The below process kinda breaks down if it *is* 0 anyways.
	if i.Sign() == 0 {
		return "0"
	}

	var str string

	//Don't want to modify original.
	z := new(big.Int).Set(i)

	//Save the sign. true if negative. false otherwise.
	negative := z.Sign() == -1
	//Absolute value to make things straight-forward.
	z.Abs(z)

	for m := new(big.Int); z.Cmp(zero) != 0; {
		// DivMod sets `z` to the division result, and `m` to the "remainder" of the operation.
		z.DivMod(z, dozen, m)
		str = Digits[int(m.Int64())] + str
	}

	if negative {
		str = "-" + str
	}
	return str
}

//Accepts a big.Rat from the `math/big` package, and returns a dozenal number as a string.
// If the number can be represented as a whole number, the result is returned. Non-whole numbers are
// represented by the whole number and fractional place-values delimetered by a `;`. If fractional
// place-values are repeating, the repeating digits are surrounded in brackets.
//
// For example, a big.Rat representing (in Decimal) 1/12 will return as "0;1", and 1/9 as "0;14".
// 1/10 should appear as "0;1[2497]", as the "2497" would normally repeat infinitely in dozenal.
func Rat(i *big.Rat) string {
	//If the Denominator is 1, just return the Numerator.
	if i.IsInt() {
		return Int(i.Num())
	}
	rem := new(big.Int)
	q := new(big.Int)
	q.DivMod(i.Num(), i.Denom(), rem)

	if rem.Cmp(big.NewInt(0)) == 0 {
		//If Num/Denom is an integer (no remainder) return the result.
		return Int(q)
	} else {
		//Otherwise, traverse the remainder, until you find a loop.
		var pseq []*big.Int //Sequence of place values.
		var rseq []*big.Int //Sequence of remainders (to know when a repeating patern begins)
		var repeat int

		//Multiply the remainder by a dozen, and then divide it by the denominator
		//Repeat until there is no remainder (0), or the remainder repeats.
		for {
			rem.Mul(rem, dozen)

			d, rem := new(big.Int).DivMod(rem, i.Denom(), rem)
			pseq = append(pseq, new(big.Int).Set(d))
			rseq = append(rseq, new(big.Int).Set(rem))
			if rem.Cmp(zero) == 0 {
				//It terminates!
				repeat = -1
				break
			} else if repeat = repeats(pseq, rseq); repeat > -1 {
				//It repeats!
				// TODO
				// This check for repitions is probably very costly.
				// Need a better way to check for repitions.
				// Possibly represent the place-values as a map[big.Int]big.Int,
				//  With the remainder as the key? Not sure if it'll work...

				//Remove the last elements, because it's a repition of the first digit in the loop.
				pseq = pseq[:len(pseq)-1]
				rseq = rseq[:len(rseq)-1]
				break
				//	} else {
				//		//It keeps going!
				//		if len(pseq) >= 10 { //Only up to 10 (for testing)
				//			repeat = -1
				//			break
				//		}
			}
		}

		var frac string
		for i, v := range pseq {
			if i == repeat {
				frac += "["
			}
			frac += Int(v)
		}
		if repeat > -1 {
			frac += "]"
		}

		return Int(q) + ";" + frac
	}
}

//Treats each iteration of ps and rems as pairs, and checks if the latest iteration is repeated.
func repeats(ps, rems []*big.Int) int {
	l := len(ps) - 1 //index of last element
	for i := range ps {
		if i == l {
			break
		}
		if ps[l].Cmp(ps[i]) == 0 && rems[l].Cmp(rems[i]) == 0 {
			return i
		}
	}
	return -1
}
