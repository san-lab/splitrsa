package main

import "fmt"

/*
Black 	30	40
Red 31	41
Green	32	42
Yellow	33	43
Blue	34	44
Magenta	35	45
Cyan	36	46
White	37	47

Option	Code	Description
Reset
0	Back to normal (remove all styles)
Bold
1	Bold the text
Underline
4	Underline text
Inverse
7	Interchange colors of background and foreground
Bold off
21	Normal from bold
Underline off
24	Normal from Underline
Inverse off
27	Reverse of the Inverse
*/
var Black, Red, Green, Yellow, Blue, Magenta, Cyan, White string
var Colors = []*string{&Black, &Red, &Green, &Yellow, &Blue, &Magenta, &Cyan, &White}

var Reset, Bold, Dim, Script, Underline, X5, X6, Swap string
var Specials = []*string{&Reset, &Bold, &Dim, &Script, &Underline, &X5, &X6, &Swap}

const SpeciallOff = 20

const font = 30
const backg = 40
const esc = "\033["
const reset = esc + "0m"

func CFormat(msg string, color, special *string) (cmsg string) {
	spec := ""
	for i := 0; i < len(Specials); i++ {
		if Specials[i] == special {
			spec = fmt.Sprint(i)
		}
	}
	cmsg = msg
	for i, c := range Colors {
		if c == color {
			ci := font + i
			cmsg = esc + spec + ";" + fmt.Sprintf("%vm", ci) + cmsg + reset
		}
	}

	return
}

// Red
func RedIt(om string, a ...any) string {

	return CFormat(fmt.Sprintf(om, a...), &Red, nil)
}

// Red bold
func RedBIt(om string, a ...any) string {
	return CFormat(fmt.Sprintf(om, a...), &Red, &Bold)
}

// Green
func GreenIt(om string, a ...any) string {

	return CFormat(fmt.Sprintf(om, a...), &Green, nil)
}

// Green bold
func GreenBIt(om string, a ...any) string {
	return CFormat(fmt.Sprintf(om, a...), &Green, &Bold)
}
