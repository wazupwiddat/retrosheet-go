package readers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/wazupwiddat/retrosheet/models"
)

const FileMatcher = "TEAM"

func ReadTeams(r *zip.ReadCloser) []models.Team {
	teams := []models.Team{}
	for _, f := range r.File {
		if !strings.Contains(f.Name, FileMatcher) {
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

		year := ParseYear(f.Name)
		for _, r := range records {
			if len(r) == 0 {
				fmt.Println("Empty record")
				continue
			}
			team := models.Team{
				Year: year,
			}
			for i, f := range r {
				switch i {
				case 0:
					team.TeamCode = f
				case 1:
					team.League = ParseLeague(f)
				case 2:
					team.Name = f
				case 3:
					team.Mascot = f
				}
			}
			teams = append(teams, team)
		}
		rc.Close()
	}
	return teams
}
