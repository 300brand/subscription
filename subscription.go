package main

import (
	"fmt"
	"github.com/300brand/subscription/samplesite"
)

func main() {
	fmt.Println(samplesite.Start())
	select {}
}
