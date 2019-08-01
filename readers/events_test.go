package readers_test

import (
	"os"
	"testing"

	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	convey "github.com/smartystreets/goconvey/convey"
	"github.com/wazupwiddat/retrosheet/readers"
)

func TestReadEventFile(t *testing.T) {
	convey.Convey("Given an event file ...", t, func() {
		r, err := os.Open("../testdata/2018ANA.EVA")
		convey.So(err, convey.ShouldBeNil)

		db, mock, err := sqlmock.New()
		convey.So(err, convey.ShouldBeNil)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM players?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "player_id", "lastname", "firstname", "1", "1"}))

		conn := &dbr.Connection{
			DB:            db,
			EventReceiver: &dbr.NullEventReceiver{},
			Dialect:       dialect.MySQL,
		}

		sess := conn.NewSession(nil)

		readers.ReadGameEventsFromFile(sess, r)

	})
}
