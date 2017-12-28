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
	read_size, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(read_size, "bytes read.", string(buf))
	fmt.Printf("%x\n", buf)
}
