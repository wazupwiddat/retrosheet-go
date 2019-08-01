package readers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"log"
	"path/filepath"

	"github.com/wazupwiddat/retrosheet/models"
)

const RosterFileExt = ".ROS"

func ReadPlayers(r *zip.ReadCloser) []models.Player {
	players := []models.Player{}
	for _, f := range r.File {
		if filepath.Ext(f.Name) != RosterFileExt {
			continue
		}
		fmt.Printf("Reading %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err, f.Name)
		}

		reader := csv.NewReader(rc)
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range records {
			if len(r) == 0 {
				fmt.Println("Empty record")
				continue
			}
			player := models.Player{}
			for i, f := range r {
				switch i {
				case 0:
					player.PlayerID = f
				case 1:
					player.LastName = f
				case 2:
					player.FirstName = f
				case 3:
					player.Bats = ParseHanded(f)
				case 4:
					player.Throws = ParseHanded(f)
				}
			}
			players = append(players, player)
		}
		rc.Close()
	}
	return players
}
