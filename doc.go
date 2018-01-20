//	This package implements a formatter useful for displaying numbers of arbitrary bases.
//
//	Usage
//
//	A formatter has methods for outputting a number to a string. The base/dozenal package
//	has some default formatters for outputting numbers in base-12, otherwise known as dozenal.
//		var num int64 = 123
//		fmt.Printf("Decimal: %d\nDozenal: %s", num, dozenal.ASCII.Int64(num))
//	Output:
//		Decimal: 123
//		Dozenal: X3
//	Notice that when using fmt, you need to use the "s" verb, because the formmater
//	already gives you a string as output.
//
//	New Formatter
//
//	If you want different numerals, such as if you have your own digits for ten and eleven
//	in Dozenal or even if you want to map 0-9 to something unique as well, you can use NewFormatter()
//	to create a formatter.
//		mybase := base.NewFormatter([]string{13, "a","b","c","d","e","f","g","h","i","j","k","l","m","n"})
//		fmt.Println(mybase.Int64(100)) // Output: "hj"
//	Notice in this example that we chose to use 13 as a base, and yet the slice literal contains more than
//	13 elements. In such 	a case, the NewFormatter function will truncate the slice to the length of the base,
//	in this case ending at "m".
//
//	Big
//
//	The formatter also has methods for the arbitrarily large numbers implemented in math/big.
//	BigInt() and BigRat() accept a big.Int and big.Rat respectively, calculate it, and output it
//	in the desired base. This allows displaying arbitrarily large integers and fractions.
//
//	big.Rat.String() outputs a number in the "a/b" format, however the formatter goes to the effort
//	of evaluating that expression. Thus, the following example:
//		r := big.NewRat(1,7) // Represents a rational number of 1/7
//		fmt.Println(r)
//		fmt.Println(dozenal.ASCII.BigRat(r))
//	Output:
//		1/7
//		0;[186X35]
//	Because big.Rat represents a rational number, the fractional form will eventually terminate or loop.
//	In this case, 1/7 is an infinite loop of 6 repeating digits, represented within square-brackets.
package base
