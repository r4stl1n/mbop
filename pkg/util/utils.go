package util

import (
	"fmt"
	"unicode"
)

type Utils struct {
}

func (u *Utils) Init() *Utils {

	*u = Utils{}

	return u
}

func (u *Utils) EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	length := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	// If here, string is shorter or equal to maxLen
	return text
}

func (u *Utils) PrintBanner() {
	fmt.Println(`  __  __ ____   ____  _____  
 |  \/  |  _ \ / __ \|  __ \ 
 | \  / | |_) | |  | | |__) |
 | |\/| |  _ <| |  | |  ___/ 
 | |  | | |_) | |__| | |     
 |_|  |_|____/ \____/|_|
-- Merry Band of Pirates - Multi Agent Automation`)
}
