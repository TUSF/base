package base

import (
	"math/big"
)

//Because the "math" package doesn't have a const for MaxInt.
//len(slice) returns an int, and slice[index] only accepts an int as an index.
const maxInt = int((^uint(0)) >> 1)

//	Formatter can be created to utilize sets of digits and bases that are atypical.
type Formatter struct {
	//The base to format numbers in.
	base  int
	bbase *big.Int

	//Not exported, because changing it isn't thread-safe.
	digits []string
}

//	Create a Formatter with the desired digits to use.
//	The base is decided by the len() of slice.
func NewFormatter(slice []string) Formatter {
	// Copy the slice.
	// Might cause a panic if someone decides to change the original slice while also using it.
	b := len(slice)
	n := make([]string, b)
	copy(n, slice)
	return Formatter{base: b, bbase: big.NewInt(int64(b)), digits: n}
}

//	Accepts an int64, and returns a number as a string.
func (z Formatter) Int64(n int64) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	var str string

	//This process is simple. Divide the integer by the base, and write the
	// remainder to the string... Well, that's technically what's going on.
	for ; n != 0; n /= int64(z.base) {
		rem := n % int64(z.base)
		str = z.digits[rem] + str
	}

	if negative {
		str = "-" + str
	}
	return str
}

//	Accepts a uint64, and returns a number as a string.
func (z Formatter) UInt64(n uint64) string {
	if n == 0 {
		return "0"
	}
	var str string

	//This process is simple. Divide the integer by the base, and write the
	// remainder to the string... Well, that's technically what's going on.
	for ; n != 0; n /= uint64(z.base) {
		rem := n % uint64(z.base)
		str = z.digits[rem] + str
	}
	return str
}

//	Accepts a big.Int from the `math/big` package, and returns an integer as a string.
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

	for m := new(big.Int); d.Sign() != 0; {
		// DivMod sets `d` to the division result, and `m` to the "remainder" of the operation.
		d.DivMod(d, z.bbase, m)
		str = z.digits[int(m.Int64())] + str
	}

	if negative {
		str = "-" + str
	}
	return str
}

//	Accepts a big.Rat from the `math/big` package, and returns a number as a string.
//	If the number can be represented as a whole number, the result is returned. Non-whole numbers are
//	represented by the whole number and fractional place-values delimetered by a `;`. If fractional
//	place-values are repeating, the repeating digits are surrounded in brackets.
//
//	For example, using base 12, a big.Rat representing (in Decimal) 1/12 will return as "0;1", and 1/9 as "0;14".
//	1/10 should appear as "0;1[2497]", as the "2497" would normally repeat infinitely in dozenal.
func (z Formatter) BigRat(i *big.Rat) string {
	//If the Denominator is 1, just return the Numerator.
	if i.IsInt() {
		return z.BigInt(i.Num())
	}
	//Gonna be reusing this `rem` for remainder a lot.
	rem := new(big.Int)
	//`q` for Quotent is basically everything before the fractional point.
	q := new(big.Int)

	//Divide the Numerator and Denominator; the quotent goes into `q`, remainder into `rem`.
	q.DivMod(i.Num(), i.Denom(), rem)

	if rem.Sign() == 0 {
		//If Num/Denom is an integer (no remainder) return the result.
		return z.BigInt(q)
	}

	//Otherwise, traverse the remainder, until you find a loop.
	var pseq []*big.Int          //Sequence of place values.
	rseq := make(map[string]int) //Sequence of remainders (to know when a repeating patern begins)
	var repeat = -1

	//Multiply the remainder by the base, and then divide it by the denominator
	//Repeat until there is no remainder (0), or the remainder repeats.
	//Stop upon hitting maxInt, so that it doesn't attenpt to go past it.
	for it := 0; it < maxInt; it++ {
		// Multiply the last remainder by the base.
		rem.Mul(rem, z.bbase)

		// Divide the product by the Denminator.
		d, rem := new(big.Int).DivMod(rem, i.Denom(), rem)

		if rem.Sign() == 0 {
			//If the remainder is 0, that means we've reached the end of the answer.
			pseq = append(pseq, new(big.Int).Set(d))
			break
		} else if r, ok := rseq[rem.String()]; ok {
			//If the remainder has shown up in a previous division,
			// then that means the sequence repeats.
			//Be sure to mark when the first occurance of this remainder occurs.
			repeat = r
			break
		} else {
			//Otherwise, mark down the digit to be written to the string later.
			pseq = append(pseq, new(big.Int).Set(d))

			//And keep track of the remainder, and its iteration.
			rseq[rem.String()] = it
		}
	}

	var frac string
	for it, v := range pseq {
		if it == repeat {
			frac += "["
		}
		frac += z.BigInt(v)
	}
	if repeat > -1 {
		frac += "]"
	}

	return z.BigInt(q) + ";" + frac

}
