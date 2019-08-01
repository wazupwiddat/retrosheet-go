package mysql

import (
	"archive/zip"
	"fmt"
	"log"
	"sync"

	"github.com/wazupwiddat/retrosheet/db"
	"github.com/wazupwiddat/retrosheet/models"
	"github.com/wazupwiddat/retrosheet/readers"
)

var lock sync.Mutex

func unique(slice []models.Player) []models.Player {
	keys := make(map[string]bool)
	list := []models.Player{}
	for _, entry := range slice {
		if _, value := keys[entry.PlayerID]; !value {
			keys[entry.PlayerID] = true
			list = append(list, entry)
		}
	}
	return list
}

func LoadTeams(r *zip.ReadCloser) error {
	conn, err := db.Open("mysql", "")
	if err != nil {
		log.Println(err)
		return err
	}
	session := conn.NewSession(nil)

	teams := readers.ReadTeams(r)
	tx, err := session.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.RollbackUnlessCommitted()

	err = models.SaveTeams(tx, teams)
	if err != nil {
		log.Println(err)
		return err
	}

	tx.Commit()
	return nil
}

func LoadPlayers(r *zip.ReadCloser) error {
	lock.Lock()

	defer lock.Unlock()

	conn, err := db.Open("mysql", "")
	if err != nil {
		log.Println(err)
		return err
	}
	session := conn.NewSession(nil)

	players := readers.ReadPlayers(r)
	players = unique(players)
	tx, err := session.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.RollbackUnlessCommitted()

	err = models.SavePlayers(tx, players)
	if err != nil {
		log.Println(err)
		return err
	}

	tx.Commit()
	return nil
}

func LoadGames(r *zip.ReadCloser) error {
	conn, err := db.Open("mysql", "")
	if err != nil {
		log.Println(err)
		return err
	}
	session := conn.NewSession(nil)

	games := readers.ReadGames(session, r)
	tx, err := session.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.RollbackUnlessCommitted()

	err = models.SaveGames(tx, games)
	if err != nil {
		log.Println(err)
		return err
	}

	tx.Commit()
	return nil
}

func LoadGamesEvents(r *zip.ReadCloser) error {
	conn, err := db.Open("mysql", "")
	if err != nil {
		log.Println(err)
		return err
	}
	session := conn.NewSession(nil)

	log.Println("Reading games events...")
	events := readers.ReadGamesEvents(session, r)
	fmt.Println(len(events))
	tx, err := session.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.RollbackUnlessCommitted()

	log.Println("Saving games events...")
	err = models.SaveGamesEvents(tx, events)
	if err != nil {
		log.Println(err)
		return err
	}

	tx.Commit()
	return nil
}
