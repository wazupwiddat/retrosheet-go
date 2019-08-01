package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wazupwiddat/retrosheet/mysql"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var eventDirectory string
	flag.StringVar(&eventDirectory, "output", "./../../output", "Read output path. Default: 'output'")
	flag.Parse()

	filename, err := filepath.Abs(filepath.Dir(fmt.Sprintf("%s/*.zip", eventDirectory)))
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(filename)
	if err != nil {
		log.Fatal(err)
	}

	g := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(g)

	archiveChannel := make(chan string, g)
	for i := 0; i < g; i++ {
		go func() {
			defer wg.Done()

			for {
				filename, ok := <-archiveChannel
				if !ok {
					log.Println("Thread Done.")
					break
				}
				r, err := zip.OpenReader(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer r.Close()

				mysql.LoadTeams(r)
				mysql.LoadPlayers(r)
				mysql.LoadGames(r)

				// game events last
				mysql.LoadGamesEvents(r)
			}
		}()
	}

	for _, f := range files {
		if !f.Mode().IsRegular() || filepath.Ext(f.Name()) != ".zip" {
			continue

		}
		filename, err := filepath.Abs(fmt.Sprintf("%s/%s", eventDirectory, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		archiveChannel <- filename
	}

	close(archiveChannel)
	wg.Wait()
}
