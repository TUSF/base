//	This package implements a formatter useful for displaying numbers as dozenal,
//	otherwise called duodecimal.
//
//	There are three default formatters to use, which each follow conventions used by
//	dozenal advocates.
//		var ASCII base.Formatter
//		var Amer base.Formatter
//		var Brit base.Formatter
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
package dozenal
