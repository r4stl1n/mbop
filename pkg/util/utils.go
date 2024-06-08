package util

import "fmt"

type Utils struct {
}

func (u *Utils) Init() *Utils {

	*u = Utils{}

	return u
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
