package recovering

import "fmt"

func ILoveCrashing() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
			panic("and this is where it ends")
		}
	}()

	panic("this is where it starts")
}
