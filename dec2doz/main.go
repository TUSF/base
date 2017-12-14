package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/TUSF/doz"
)

func main() {
	doz.Digits = doz.AmerDigits
	s := bufio.NewScanner(os.Stdin)
	var INT big.Int
	var RAT big.Rat

	for s.Scan() {
		b := s.Text()
		if b == "quit" || b == "q" {
			return
		}

		if INT, t := INT.SetString(b, 0); t {
			//First, assume it's a plain integer.
			fmt.Println(doz.Int(INT))

		} else if RAT, t := RAT.SetString(b); t {
			//Second, assume it's a fraction.
			fmt.Println(doz.Rat(RAT))

		} else {
			//Third, assume it's a decimal number.
			if strings.Index(b, ".") > -1 {
				if nums := strings.Split(b, "."); len(nums) == 2 {
					//Convert each side of the decimal point into a big.Int?
					if nums[1] == "" {
						fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
						continue
					}
					if nums[0] == "" {
						RAT.SetInt64(0)
					} else {
						RAT, t := RAT.SetString(nums[0])
						if !t {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}
						if _, err := strconv.Atoi(nums[1]); err != nil {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
							continue
						}

						ten := big.NewInt(10)
						for i := 0; i < len(nums[1]); i++ {
							d, _ := strconv.Atoi(string(nums[1][i]))
							dn := big.NewInt(int64(d))
							in := new(big.Int).SetInt64(int64(i + 1))
							RAT.Add(RAT, new(big.Rat).SetFrac(dn, new(big.Int).Exp(ten, in, nil)))
							// 1.123 = 1 + 1/10 + 2/100 + 3/1000
							// Of course, once you convert it into dozenal, even a simple decimal like that becomes huge.
							// "1.123" becomes 1;15[…], followed by an infinitely repeating sequence of 50+ digits
						}
						fmt.Println(doz.Rat(RAT))
					}
				} else {
					fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
				}
			} else {
				fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
			}
		}

		/*
			//First, assume it's a plain integer.
			INT, t := INT.SetString(b, 0)
			if !t {
				//Second, assume it's a fraction.
				RAT, t := RAT.SetString(b)
				if !t {
					//Third, assume it's a decimal number.
					if strings.Index(b, ".") > -1 {
						if nums := strings.Split(b, "."); len(nums) == 2 {
							//Convert each side of the decimal point into a big.Int?
						} else {
							fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
						}
					} else {
						fmt.Fprintln(os.Stderr, "Not a valid number. Integers, Fractions or Decimals only.")
						continue
					}
				} else {
					fmt.Println(doz.Rat(RAT))
				}
			} else {
				fmt.Println(doz.Int(INT))
			}
		*/
	}
}
