package internal

import "fmt"

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
