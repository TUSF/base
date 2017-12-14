package dozenal

import (
	"math/big"
)

//Likely to be initialized a lot, so may as well reuse the same one.
var zero *big.Int = big.NewInt(0)
var dozen *big.Int = big.NewInt(0xC)

//	The ASCII formatter uses an ASCII-friendly set of Dozenal digits to be compatible with most displays.
var ASCII Formatter = NewFormatter([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "X", "E"})

//	The Amer formatter uses digits recommended by The Dozenal Society of America. Dec and El are
//	represented using the Greek letter, Chi (Χ) and the Latin "open E". (Ɛ)
var Amer Formatter = NewFormatter([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u03A7", "\u0190"})

//	The Brit formatter uses digits recommended by The Dozenal Society of Great Britain . Dec and El are
//	represented using a turned 2 (↊) and 3 (↋). These are the least supported among fonts.
var Brit Formatter = NewFormatter([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u218A", "\u218B"})

//	A formatter can be created to utilize sets of dozenal digits that are atypical. Under numeral circumstances
//	simply using one of the default Formatters should be adequate.
type Formatter struct {
	//Not exported, because changing it isn't thread-safe.
	digits []string
}

//	Create a Formatter with the desired digits to use. Slices longer than a dozen elements will be truncated.
func NewFormatter(slice []string) Formatter {
	// Copy the slice, but only the first 12 elements.
	// Might cause a panic if someone decides to change the original slice while also using it.
	n := make([]string, 0xC)
	copy(n, slice)
	return Formatter{digits: n}
}

//	Accepts an int64, and returns a dozenal number as a string.
func (z Formatter) Int64(n int64) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	var str string

	//This process is simple. Divide the integer by a dozen (0xC), and write the
	// remainder to the string... Well, that's technically what's going on.
	for ; n != 0; n /= 0xC {
		rem := n % 0xC
		str = z.digits[rem] + str
	}

	if negative {
		str = "-" + str
	}
	return str
}

//	Accepts a uint64, and returns a dozenal number as a string.
func (z Formatter) UInt64(n uint64) string {
	if n == 0 {
		return "0"
	}
	var str string

	//This process is simple. Divide the integer by a dozen (0xC), and write the
	// remainder to the string... Well, that's technically what's going on.
	for ; n != 0; n /= 0xC {
		rem := n % 0xC
		str = z.digits[rem] + str
	}
	return str
}

//	Accepts a big.Int from the `math/big` package, and returns a dozenal integer as a string.
func (z Formatter) BigInt(i *big.Int) string {
	//No need to waste time if it's 0. The below process kinda breaks down if it *is* 0 anyways.
	if i.Sign() == 0 {
		return "0"
	}

	var str string

	//Don't want to modify original.
	d := new(big.Int).Set(i)

	//Save the sign. true if negative. false otherwise.
	negative := d.Sign() == -1
	//Absolute value to make things straight-forward.
	d.Abs(d)

	for m := new(big.Int); d.Cmp(zero) != 0; {
		// DivMod sets `d` to the division result, and `m` to the "remainder" of the operation.
		d.DivMod(d, dozen, m)
		str = z.digits[int(m.Int64())] + str
	}

	if negative {
		str = "-" + str
	}
	return str
}

//	Accepts a big.Rat from the `math/big` package, and returns a dozenal number as a string.
//	If the number can be represented as a whole number, the result is returned. Non-whole numbers are
//	represented by the whole number and fractional place-values delimetered by a `;`. If fractional
//	place-values are repeating, the repeating digits are surrounded in brackets.
//
//	For example, a big.Rat representing (in Decimal) 1/12 will return as "0;1", and 1/9 as "0;14".
//	1/10 should appear as "0;1[2497]", as the "2497" would normally repeat infinitely in dozenal.
func (z Formatter) BigRat(i *big.Rat) string {
	//If the Denominator is 1, just return the Numerator.
	if i.IsInt() {
		return z.BigInt(i.Num())
	}
	rem := new(big.Int)
	q := new(big.Int)
	q.DivMod(i.Num(), i.Denom(), rem)

	if rem.Cmp(zero) == 0 {
		//If Num/Denom is an integer (no remainder) return the result.
		return z.BigInt(q)
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
				//} else if repeat = z.repeats(pseq, rseq); repeat > -1 {
			} else if repeat = z.repeats(rseq); repeat > -1 {
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
			}
		}

		var frac string
		for i, v := range pseq {
			if i == repeat {
				frac += "["
			}
			frac += z.BigInt(v)
		}
		if repeat > -1 {
			frac += "]"
		}

		return z.BigInt(q) + ";" + frac
	}
}

//	Treats each iteration of ps and rems as pairs, and checks if the latest iteration is repeated.
func (z Formatter) repeats(rems []*big.Int) int {
	l := len(rems) - 1 //index of last element
	for i := range rems {
		if i == l {
			break
		}
		if rems[l].Cmp(rems[i]) == 0 {
			return i
		}
	}
	return -1
}
