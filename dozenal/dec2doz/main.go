//	Take stdin as a number, output it as dozenal.
package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/TUSF/base/dozenal"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		var INT big.Int
		var RAT big.Rat
		b := strings.TrimSpace(s.Text())
		if b == "quit" || b == "q" {
			return
		}

		if INT, t := INT.SetString(b, 0); t {
			//First, assume it's a plain integer.
			fmt.Println(dozenal.Amer.BigInt(INT))

		} else if RAT, t := RAT.SetString(b); t {
			//Second, assume it's a fraction. ("12/7")
			fmt.Println(dozenal.Amer.BigRat(RAT))

		} else {
			//Third, assume it's a decimal number. ("10.123")
			if strings.Index(b, ".") > -1 {
				if nums := strings.Split(b, "."); len(nums) == 2 {
					//Convert each side of the decimal point into a big.Int?
					if nums[1] == "" {
						// Nothing after the point. "10." is treated as just "10"
						if INT, t := INT.SetString(b, 0); t {
							fmt.Println(dozenal.Amer.BigInt(INT))
						} else {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
						}
						continue
					} else {
						if nums[0] == "" {
							// ".124" is treated as "0.124"
							nums[0] = "0"

						} else if _, err := strconv.Atoi(nums[0]); err != nil {
							// First make sure each side of the point is actually a decimal number.
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}
						if _, err := strconv.Atoi(nums[1]); err != nil {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}

						// 1.123 = 1 + 123/1000
						// So, treat everything before the point as an integer
						// then feed the number after the point, divided by the next power of 10.
						//
						// Of course, once you convert it into dozenal, even a simple decimal like that becomes huge.
						// "1.123" becomes 1;15[…], followed by an infinitely repeating sequence of 50+ digits
						RAT, t := RAT.SetString(nums[0])
						if !t {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}

						d, t := new(big.Rat).SetString(
							fmt.Sprintf("%s/1%s", nums[1], strings.Repeat("0", len(nums))))
						if !t {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}
						RAT.Add(RAT, d)
						fmt.Println(dozenal.Amer.BigRat(RAT))
					}
				} else {
					// More than 1 point!?
					fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
				}
			} else {
				// Missing a point.
				fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
			}
		}
	}
}
