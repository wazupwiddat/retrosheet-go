package main

import (
	"flag"
	"log"
	"os"

	"database/sql"

	"github.com/pressly/goose"
	_ "github.com/wazupwiddat/retrosheet/cmd/schema/migrations"

	// Init DB drivers.
	_ "github.com/go-sql-driver/mysql"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbstring, command := args[0], args[1]

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
		return
	}

	db, err := sql.Open("mysql", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
		return
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}

}
