package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	HomePath = os.Getenv("HOME")
	aacDir   = filepath.Join(HomePath, "aac319480727")
)

func main() {
	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, 1)
		filename = "bananamoon"
	)

	files, err := ioutil.ReadDir(aacDir)
	if err != nil {
		log.Fatal(err)
	}

	tempdir, err := TempDiraac()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	for _, f := range files {
		wg.Add(1)
		go func(fname, tempdir string) {
			defer wg.Done()

			buf, err := AACtoByte(fname)
			if err != nil {
				errChan <- err
			}

			c := 0
			for i, _ := range buf {
				if fmt.Sprintf("%x", buf[i]) == "5c" && fmt.Sprintf("%x", buf[i+1]) == "ff" {
					c = i + 1
					break
				}
			}

			if err := createAAC(fmt.Sprintf("%s/%s", tempdir, fname), buf[c:]); err != nil {
				errChan <- err
			}
		}(f.Name(), tempdir)
	}
	select {
	case err := <-errChan:
		log.Fatal(err)
	default:
	}
	wg.Wait()

	if err := RemakeAAC(filename, tempdir); err != nil {
		log.Fatal(err)
	}
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

func TempDiraac() (string, error) {
	aacDir, err := ioutil.TempDir(HomePath, "output")
	if err != nil {
		return "", err
	}
	return aacDir, nil
}

func AACtoByte(fname string) ([]byte, error) {
	fpath := filepath.Join(aacDir, fname)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	f, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, f.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func RemakeAAC(filename, fpath string) error {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(fpath)
	if err != nil {
		return err
	}

	var res []byte
	for _, f := range files {
		res = append(res, f.Name()...)
		res = append(res, '|')
	}
	name := fmt.Sprintf("concat:%s", string(res[:len(res)-1]))

	aacFile := filepath.Join(HomePath, "Output", filename)
	cmd := exec.Command(cmdPath, "-i", name, "-c", "copy", fmt.Sprintf("%s.aac", aacFile))
	cmd.Dir = fpath
	cmd.Run()

	return nil
}
