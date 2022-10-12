package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/profile"
)

var (
	b       = make([]byte, 32)
	bufSize = 1024 * len(b)
)

func main() {
	defer profile.Start(profile.MemProfileRate(1), profile.MemProfile, profile.ProfilePath(".")).Stop()
	filename := os.Args[1]
	store := os.Args[2]
	if store == "" || filename == "" {
		fmt.Printf("filename and store must be specified")
		os.Exit(1)
	}
	defer recovery()

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fbuff := bufio.NewReaderSize(f, bufSize)

	s, err := os.OpenFile(store, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer s.Close()
	sbuff := bufio.NewWriterSize(s, bufSize)

	var wg sync.WaitGroup
	wg.Add(1)
	go write(fbuff, sbuff, &wg)
	wg.Wait()
}

func recovery() {
	if err := recover(); err != nil {
		fmt.Printf("error while processing file: %s\n", err)
	}
}

func write(r io.Reader, w io.Writer, wg *sync.WaitGroup) {
	var err error
	defer wg.Done()
	for {
		if _, err = r.Read(b[:]); err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if _, err = w.Write(b); err != nil {
			panic(err)
		}
	}
}
