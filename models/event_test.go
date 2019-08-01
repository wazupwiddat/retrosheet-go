package models_test

import (
	"testing"

	"github.com/wazupwiddat/retrosheet/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
	"github.com/smartystreets/goconvey/convey"
)

func TestEventSave(t *testing.T) {
	convey.Convey("Given a Game event  ...", t, func() {
		db, mock, err := sqlmock.New()
		convey.So(err, convey.ShouldBeNil)
		defer db.Close()

		mock.ExpectExec("INSERT INTO `game_events` (.+) VALUES (.+)").
			WillReturnResult(sqlmock.NewResult(1, 1))
		conn := &dbr.Connection{
			DB:            db,
			EventReceiver: &dbr.NullEventReceiver{},
			Dialect:       dialect.MySQL,
		}

		sess := conn.NewSession(nil)

		ge := models.NewGameEvent(1, 4, 0, 0, 1)
		ge.Play = models.EventDetail{
			Play:       models.GroundRuleDouble,
			ExtraPlays: []models.BasicPlay{},
			Fielders:   []models.Position{},
			Modifiers: []models.Modifier{
				{
					PlayModifier: models.ModifierLinedDrive,
					Location:     "9LS",
				},
			},
			RunnerAdv: []models.RunnerAdvance{
				{
					StartBase:  2,
					FinishBase: 4,
				},
			},
		}
		err = ge.Save(sess)
		convey.So(err, convey.ShouldBeNil)
	})
}
