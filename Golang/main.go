package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Since(time.UnixMilli(1751179796000)))
}
