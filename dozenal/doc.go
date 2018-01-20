//	This package implements a formatter useful for displaying numbers as dozenal,
//	otherwise called duodecimal.
//
//	There are three default formatters to use, which each follow conventions used by
//	dozenal advocates.
//		var ASCII Formatter
//		var Amer Formatter
//		var Brit Formatter
//	The ASCII formatter uses the Latin letters X and E for the digits representing
//	ten and eleven, making it compatible with any display that already supports ASCII.
//
//	The Amer formatter uses the Greek letter, Chi (Χ) and the Latin "open E" (Ɛ) to
//	represent ten and eleven. This should work on systems with fonts for Greek and Extended
//	Latin, so it may be compatible with most modern computers. This notation is the one used
//	by The Dozenal Society of America, hence the shortened name.
//
//	The Brit formatter uses ↊ and ↋, which are a turned "2" and "3" symbol. These symbols
//	were only added to the Unicode standard in version 8 (2015), and due to their low
//	priority, they are not likely to be included on a computer's default fonts. This notation
//	is the one used by The Dozenal Society of Great Britain.
//
//	Usage
//
//	A formatter has methods for outputting a number to a string.
//		var num int64 = 123
//		fmt.Printf("Decimal: %d\nDozenal: %s", num, dozenal.ASCII.Int64(num))
//	Output:
//		Decimal: 123
//		Dozenal: X3
//	Notice that when using fmt, you need to use the "s" verb, because the dozenal formmater
//	already gives you a string as output.
//
//	New Formatter
//
//	If you want different numerals, such as if you have your own digits for ten and eleven
//	or even if you want to map 0-9 to something unique as well, you can use NewFormatter()
//	to create a formatter.
//		mydoz := dozenal.NewFormatter([]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n"})
//		fmt.Println(mydoz.Int64(100)) // Output: "ie"
//	Notice in this example that the slice literal contains more than a dozen elements. In such
//	a case, the NewFormatter function will truncate the slice to just a dozen, in this case ending at "l".
//
//	Big
//
//	The formatter also has methods for the arbitrarily large numbers implemented in math/big.
//	BigInt and BigRat accept a big.Int and big.Rat respectively, calculate it, and output it
//	in dozenal notation. This allows displaying arbitrarily large dozenal integers and fractions.
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
package dozenal
