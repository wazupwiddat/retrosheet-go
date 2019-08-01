package models

import (
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
)

type Handed int

const (
	RightHanded Handed = 1
	LeftHanded  Handed = 2
	BothHanded  Handed = 3
)

func (h Handed) String() string {
	HandedName := [...]string{
		"Right",
		"Left",
		"Both",
	}
	if h < RightHanded || h > BothHanded {
		return "Invalid Handed"
	}
	return HandedName[h]
}

type Player struct {
	ID        int
	PlayerID  string `db:"player_id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Bats      Handed
	Throws    Handed
}

func (p *Player) Save(session dbr.SessionRunner) error {
	buf := dbr.NewBuffer()

	stmt := session.InsertInto("players").
		Columns("player_id", "lastname", "firstname", "throws", "bats").
		Record(p)
	stmt.Build(dialect.MySQL, buf)

	stmt2 := session.UpdateBySql(" ON DUPLICATE KEY UPDATE lastname = ?, firstname = ?", p.LastName, p.FirstName)
	stmt2.Build(dialect.MySQL, buf)

	query, err := dbr.InterpolateForDialect(buf.String(), buf.Value(), dialect.MySQL)
	if err != nil {
		return err
	}

	_, err = session.InsertBySql(query).Exec()
	return err
}

func SavePlayers(session dbr.SessionRunner, players []Player) error {
	var err error
	for _, p := range players {
		err = p.Save(session)
		if err != nil {
			break
		}
	}
	return err
}

func GetPlayer(session dbr.SessionRunner, playerID string) (Player, error) {
	player := Player{}
	_, err := session.Select("*").From("players").
		Where("players.player_id=?", playerID).Load(&player)
	return player, err
}
