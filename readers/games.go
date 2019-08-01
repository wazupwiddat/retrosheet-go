package readers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/gocraft/dbr"
	"github.com/wazupwiddat/retrosheet/models"
)

func ReadGames(sess *dbr.Session, r *zip.ReadCloser) []models.Game {
	games := []models.Game{}
	for _, f := range r.File {
		if filepath.Ext(f.Name) != ".EVA" && filepath.Ext(f.Name) != ".EVN" {
			continue
		}
		fmt.Printf("Reading %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err, f.Name)
		}

		year := ParseYear(f.Name)
		var game models.Game
		reader := NewGameReader(rc)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				games = append(games, game)
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if len(record) > 0 {
				recordType, ok := ParseEventType(record[0])
				if !ok {
					log.Println("Missing EventType: ", record[0])
					continue
				}
				switch recordType {
				case models.GameID:
					if game.GameID != "" && game.GameID != record[1] {
						games = append(games, game)
					}
					game = models.NewGame(record[1])
				case models.Info:
					infoType, ok := ParseInfoType(record[1])
					if !ok {
						//log.Println("Missing InfoType: ", record[1])
						continue
					}
					switch infoType {
					case models.VisitingTeam:
						t, err := models.GetTeam(sess, record[2], year)
						if err != nil {
							log.Println("Failed to Find team: ", record[2], year)
							continue
						}
						game.Visitor = t.ID
					case models.HomeTeam:
						t, err := models.GetTeam(sess, record[2], year)
						if err != nil {
							log.Println("Failed to Find team: ", record[2], year)
							continue
						}
						game.Home = t.ID
					case models.GameDate:
						played, err := time.Parse("2006/01/02", record[2])
						if err != nil {
							log.Println("Failed to parse game date", err)
						}
						game.Played = played
					}
				}
			}
		}
		rc.Close()
	}
	return games
}
