package main

import (
	"fmt"
	"projecta.com/me/verification"
)

func main () {
	fmt.Println("start Auth Service")
	verification.Server()
}
