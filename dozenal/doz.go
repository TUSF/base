package dozenal

import (
	"github.com/TUSF/base"
)

//	The ASCII formatter uses an ASCII-friendly set of Dozenal digits to be compatible with most displays.
var ASCII base.Formatter = base.NewFormatter(0xC,
	[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "X", "E"})

//	The Amer formatter uses digits recommended by The Dozenal Society of America. Dec and El are
//	represented using the Greek letter, Chi (Χ) and the Latin "open E". (Ɛ)
var Amer base.Formatter = base.NewFormatter(0xC,
	[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u03A7", "\u0190"})

//	The Brit formatter uses digits recommended by The Dozenal Society of Great Britain . Dec and El are
//	represented using a turned 2 (↊) and 3 (↋). These are the least supported among fonts.
var Brit base.Formatter = base.NewFormatter(0xC,
	[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "\u218A", "\u218B"})
