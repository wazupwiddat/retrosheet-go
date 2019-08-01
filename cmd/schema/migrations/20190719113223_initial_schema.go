package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInitialSchema, downInitialSchema)
}

func upInitialSchema(txn *sql.Tx) error {
	_, err := txn.Exec(
		"CREATE TABLE `teams` (" +
			"`id` int(11) NOT NULL AUTO_INCREMENT," +
			"`team_code` varchar(4) DEFAULT NULL," +
			"`year` int(11) NOT NULL," +
			"`league` tinyint(1) NOT NULL," +
			"`name` varchar(20) DEFAULT NULL," +
			"`mascot` varchar(20) DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"KEY `team_code_idx` (`team_code`,`year`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;",
	)
	if err != nil {
		return err
	}
	_, err = txn.Exec(
		"CREATE TABLE `players` (" +
			"`id` int(11) NOT NULL AUTO_INCREMENT," +
			"`player_id` varchar(10) DEFAULT NULL UNIQUE," +
			"`lastname` varchar(30) DEFAULT NULL," +
			"`firstname` varchar(30) DEFAULT NULL," +
			"`throws` tinyint(3) NOT NULL," +
			"`bats` tinyint(3) NOT NULL," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;",
	)
	if err != nil {
		return err
	}
	_, err = txn.Exec(
		"CREATE TABLE `games` (" +
			"`id` int(11) NOT NULL AUTO_INCREMENT," +
			"`game_id` varchar(14) DEFAULT NULL," +
			"`played` datetime DEFAULT NULL," +
			"`visitor` int(11) NOT NULL," +
			"`home` int(11) NOT NULL," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;",
	)
	if err != nil {
		return err
	}
	_, err = txn.Exec(
		"CREATE TABLE `game_events` (" +
			"`id` int(11) NOT NULL AUTO_INCREMENT," +
			"`game_id` int(11) NOT NULL," +
			"`player_id` int(11) NOT NULL," +
			"`inning` int(11) NOT NULL," +
			"`inning_half` int(11) NOT NULL," +
			"`event` int(11) NOT NULL," +
			"`event_detail` TEXT," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;",
	)
	if err != nil {
		return err
	}
	return nil
}

func downInitialSchema(txn *sql.Tx) error {
	_, err := txn.Exec("DROP TABLE `teams`")
	if err != nil {
		return err
	}
	_, err = txn.Exec("DROP TABLE `players`")
	if err != nil {
		return err
	}
	_, err = txn.Exec("DROP TABLE `games`")
	if err != nil {
		return err
	}
	_, err = txn.Exec("DROP TABLE `game_events`")
	if err != nil {
		return err
	}
	return nil
}
