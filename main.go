package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("20171223_010000_g57dD.aac")
	if err != nil {
		log.Fatal(err)
	}

	f, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, f.Size())
	_, err = file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	var c int
	for i, _ := range buf {
		if fmt.Sprintf("%x", buf[i]) == "5c" {
			c = i + 1
			break
		}
	}

	a := buf[c:]
	fmt.Printf("%x", a)
}
