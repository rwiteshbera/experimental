package main

import (
	"log"
	"math"

	json "encoding/json/v2"
)

type Data struct {
	Pos float64 `json:"pos,format:nonfinite"`
}

func main() {
	d := Data{Pos: math.Inf(+1)}

	v, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(v)) // {"pos":"Infinity"}
}
