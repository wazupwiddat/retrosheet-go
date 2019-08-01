package models

import "github.com/gocraft/dbr"

type League int

const (
	American League = 0
	National League = 1
)

func (l League) String() string {
	LeagueNames := [...]string{
		"American",
		"National",
	}
	if l < American || l > National {
		return "Invalid League"
	}
	return LeagueNames[l]
}

type Team struct {
	ID       int
	TeamCode string `db:"team_code"`
	Year     int
	Name     string
	Mascot   string
	League   League
}

func NewTeam(year int, name string) Team {
	team := Team{
		Year: year,
		Name: name,
	}
	return team
}

func (t *Team) Save(session dbr.SessionRunner) error {
	_, err := session.InsertInto("teams").
		Columns("team_code", "year", "name", "mascot", "league").
		Record(t).
		Exec()
	return err
}

func SaveTeams(session dbr.SessionRunner, teams []Team) error {
	var err error
	for _, t := range teams {
		err = t.Save(session)
		if err != nil {
			break
		}
	}
	return err
}

func GetTeam(session dbr.SessionRunner, teamID string, year int) (Team, error) {
	team := Team{}
	_, err := session.Select("*").From("teams").
		Where("teams.team_code=? AND teams.year=?", teamID, year).Load(&team)
	return team, err
}
