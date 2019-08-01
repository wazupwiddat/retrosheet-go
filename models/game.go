package models

import (
	"time"

	"github.com/gocraft/dbr"
)

type Position int

const (
	PositionPitcher Position = iota + 1
	PositionCatcher
	PositionFirstBase
	PositionSecondBase
	PositionThirdBase
	PositionShortStop
	PositionLeftField
	PositionCenterField
	PositionRightField
)

func (p Position) String() string {
	PositionName := [...]string{
		"Pitcher",
		"Catcher",
		"First Base",
		"Second Base",
		"Third Base",
		"Short Stop",
		"Left Field",
		"Center Field",
		"Right Field",
	}
	if p < PositionPitcher || p > PositionRightField {
		return "Invalid Position"
	}
	return PositionName[p-1]
}

type Game struct {
	ID      int
	GameID  string `db:"game_id"`
	Visitor int
	Home    int
	Played  time.Time
}

func NewGame(gameID string) Game {
	g := Game{
		GameID: gameID,
	}
	return g
}

func (g *Game) Save(session dbr.SessionRunner) error {
	_, err := session.InsertInto("games").
		Columns("game_id", "played", "visitor", "home").
		Record(g).
		Exec()
	return err
}

func SaveGames(session dbr.SessionRunner, games []Game) error {
	var err error
	for _, g := range games {
		err = g.Save(session)
		if err != nil {
			break
		}
	}
	return err
}

func GetGame(session dbr.SessionRunner, gameID string) (Game, error) {
	game := Game{}
	_, err := session.Select("*").From("games").
		Where("games.game_id=?", gameID).Load(&game)
	return game, err
}
