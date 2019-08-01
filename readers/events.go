package readers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/gocraft/dbr"
	"github.com/wazupwiddat/retrosheet/models"
)

func ReadGamesEvents(sess *dbr.Session, r *zip.ReadCloser) []models.GameEvent {
	gameEvents := []models.GameEvent{}
	for _, f := range r.File {
		if filepath.Ext(f.Name) != ".EVA" && filepath.Ext(f.Name) != ".EVN" {
			continue
		}
		fmt.Printf("Reading %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err, f.Name)
		}

		g := ReadGameEventsFromFile(sess, rc)
		gameEvents = append(gameEvents, g...)
	}
	return gameEvents
}

func ReadGameEventsFromFile(sess *dbr.Session, file io.ReadCloser) []models.GameEvent {
	gameEvents := []models.GameEvent{}
	// year := ParseYear(f.Name)
	var game models.Game
	var gameEvent models.GameEvent
	reader := NewGameReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			gameEvents = append(gameEvents, gameEvent)
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
				game, err = models.GetGame(sess, record[1])
				if err != nil {
					log.Println("Failed to Find game: ", record[1], err)
					continue
				}
			case models.Play:
				inning := ParseInning(record[1])
				half, _ := ParseInningHalf(record[2])
				player, err := models.GetPlayer(sess, record[3])
				if err != nil {
					log.Println("Failed to Find player: ", record[3], err)
					continue
				}
				// cnt := ParseBallsStrikes(record[4])
				// pitches := ParsePitches(record[5])

				gameEvent := models.NewGameEvent(game.ID,
					models.Play, inning, half, player.ID)
				// Parse the actual events
				eventDetail, ok := ParseEventDetail(record[6])
				if !ok {
					log.Println("Parse Event Detail failed: ", record[6])
					continue
				}
				gameEvent.Play = eventDetail
				// log.Println(record, "\n", gameEvent)
				gameEvents = append(gameEvents, gameEvent)
			}
		}
	}
	return gameEvents
}
