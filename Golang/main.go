package main

import (
	"log"
	"strings"
)

func main() {
	key := "abc:cde:fgh"
	abc := strings.SplitN(key, ":", 3)
	log.Println(abc[len(abc)-1])
}
