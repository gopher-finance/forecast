package term

import "github.com/fatih/color"

/* prints either green or red text to the screen, depending
 * on decision. */
func Binary(text string, decision bool) string {
	if decision {
		return color.GreenString(text)
	}
	return color.RedString(text)
}

func Binaryf(f float64, decision bool) string {
	if decision {
		return color.GreenString("%.2f", f)
	}
	return color.RedString("%.2f", f)
}

func Binaryfp(f float64, decision bool) string {
	if decision {
		return color.GreenString("%.2f%%", f)
	}
	return color.RedString("%.2f%%", f)
}

func Arrow(decision bool) string {
	if decision {
		return "↑"
	}
	return "↓"
}
