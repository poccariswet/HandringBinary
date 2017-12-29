package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	files := []string{
		"aac/20171223_010000_g57dD.aac",
		"aac/20171223_010005_k3hXp.aac",
		"aac/20171223_010010_RkhrL.aac",
		"aac/20171223_010015_n7dqI.aac",
	}

	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go func(name string, i int) {
			defer wg.Done()
			file, err := os.Open(name)
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

			createAAC(fmt.Sprintf("output/output%d.aac", i), buf[c:])
		}(files[i], i)
	}
	wg.Wait()
}

func createAAC(name string, bf []byte) error {
	wf, err := os.Create(name)
	if err != nil {
		return err
	}
	defer wf.Close()

	wf.Write(bf)

	return nil
}
