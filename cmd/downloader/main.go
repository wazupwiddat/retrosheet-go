package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var invalidYear = [14]int{
	1923, 1924, 1925, 1926, 1928, 1929, 1930, 1932, 1933, 1934, 1935, 1936, 1937, 1939,
}

func main() {
	var startYear, endYear int
	var outputDirectory string

	currentYear := time.Now().Year()
	flag.IntVar(&startYear, "start", 1921, "Start year. Default: 1921")
	flag.IntVar(&endYear, "end", currentYear, "Start year. Default: 1921")
	flag.StringVar(&outputDirectory, "output", "output", "Download output path. Default: '.'")
	flag.Parse()

	os.MkdirAll(outputDirectory, os.ModePerm)

	workers := runtime.NumCPU()
	wch := make(chan int, workers)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for {
				year, ok := <-wch
				if !ok {
					fmt.Println("thread done")
					break
				}

				downloadYear(year, outputDirectory)
				fmt.Printf("Download Complete(%d)...\n", year)
			}
		}()
	}

	for y := startYear; y <= endYear; y++ {
		if ok := validYear(y); !ok {
			continue
		}
		wch <- y
	}
	close(wch)
	wg.Wait()
}

func validYear(year int) bool {
	for _, i := range invalidYear {
		if i == year {
			return false
		}
	}
	return true
}

func downloadYear(year int, dir string) {
	filename := fmt.Sprintf("%deve.zip", year)
	filePath := filepath.Join(dir, filename)
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	url := fmt.Sprintf("http://www.retrosheet.org/events/%deve.zip", year)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}
