package readers_test

import (
	"testing"

	"github.com/wazupwiddat/retrosheet/readers"

	convey "github.com/smartystreets/goconvey/convey"
	"github.com/wazupwiddat/retrosheet/models"
)

func TestParseRunnerAdvances(t *testing.T) {
	convey.Convey("Given runner advance string...", t, func() {
		tests := []struct {
			t []string
			v []models.RunnerAdvance
		}{
			{
				[]string{"BXH(E8/TH)(UR)"},
				[]models.RunnerAdvance{},
			},
			{
				[]string{"B-H(E8/TH)(UR)"},
				[]models.RunnerAdvance{
					{
						StartBase:  0,
						FinishBase: 4,
					},
				},
			},
			{
				[]string{"1-2#"},
				[]models.RunnerAdvance{
					{
						StartBase:  1,
						FinishBase: 2,
					},
				},
			},
		}
		for _, test := range tests {
			convey.Convey("Parsing runners"+test.t[0]+"...", func() {
				ed := readers.ParseRunnerAdvances(test.t)
				convey.So(ed, convey.ShouldResemble, test.v)
			})
		}
	})
}

func TestPlayMatching(t *testing.T) {
	convey.Convey("Given play string...", t, func() {
		tests := []struct {
			t string
			v models.BasicPlay
		}{
			{
				"SB2",
				models.StolenBase,
			},
			{
				"SBH;SB3",
				models.StolenBase,
			},
			{
				"POCSH(1361)",
				models.PickOffCaughtStealing,
			},
			{
				"PO1(16343)",
				models.PickOff,
			},
			{
				"CS2(2E4)",
				models.CaughtStealing,
			},
			{
				"CS3(23)",
				models.CaughtStealing,
			},
			{
				"CSH(12)",
				models.CaughtStealing,
			},
			{
				"IW+PO3",
				models.IntentionalWalk,
			},
			{
				"IW+SB3",
				models.IntentionalWalk,
			},
			{
				"W+WP",
				models.Walk,
			},
			{
				"W+SB3",
				models.Walk,
			},
			{
				"IW",
				models.IntentionalWalk,
			},
			{
				"W",
				models.Walk,
			},
			{
				"K",
				models.StrikeOut,
			},
			{
				"K+PB",
				models.StrikeOut,
			},
			{
				"K23+WP",
				models.StrikeOut,
			},
			{
				"K+SB2",
				models.StrikeOut,
			},
			{
				"HP",
				models.HitByPitch,
			},
			{
				"HR9",
				models.HomeRun,
			},
			{
				"H",
				models.HomeRun,
			},
			{
				"FLE5",
				models.ErrorOnFlyBall,
			},
			{
				"FC5",
				models.FieldersChoice,
			},
			{
				"E3",
				models.Error,
			},
			{
				"E1",
				models.Error,
			},
			{
				"DGR",
				models.GroundRuleDouble,
			},
			{ // T$
				"T9",
				models.Triple,
			},
			{ // D$
				"D7",
				models.Double,
			},
			{ // S$
				"S9",
				models.Single,
			},
			{ // C/E2.1-2
				"C",
				models.CatcherInterference,
			},
			{ // $(B)$(%)
				"1(B)16(2)63(1)",
				models.LinedIntoTriplePlay,
			},
			{ // $(B)$(%)
				"8(B)84(2)",
				models.LinedIntoDoublePlay,
			},
			{ // $(B)$(%)
				"3(B)3(1)",
				models.LinedIntoDoublePlay,
			},
			{ // $(%)$ $$(%)$
				"64(1)3",
				models.GroundedIntoDoublePlay,
			},
			{ // $(%)$ $$(%)$
				"6(1)3",
				models.GroundedIntoDoublePlay,
			},
			{
				"63",
				models.GroundBallOut,
			},
			{
				"8",
				models.FlyBallOut,
			},
		}
		for _, test := range tests {
			convey.Convey("Parse "+test.t+"...", func() {
				ed := readers.ParseBasicPlay(test.t)
				convey.So(ed, convey.ShouldEqual, test.v)
			})
		}
	})
}

func TestParseEventDetail(t *testing.T) {
	convey.Convey("Given play details...", t, func() {
		tests := []struct {
			t   string
			v   models.EventDetail
			err error
		}{
			{
				"DGR/L9LS.2-H",
				models.EventDetail{
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
				},
				nil,
			},
			{
				"8(B)84(2)/LDP/L8",
				models.EventDetail{
					Play:       models.LinedIntoDoublePlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionCenterField,
						models.PositionCenterField,
						models.PositionSecondBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierLinedIntoDoublePlay,
						},
						{
							PlayModifier: models.ModifierLinedDrive,
							Location:     "8",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"54(1)/FO/G5.3-H;B-1",
				models.EventDetail{
					Play:       models.GroundBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionThirdBase,
						models.PositionSecondBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierForceOut,
						},
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "5",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
						{
							StartBase:  0,
							FinishBase: 1,
						},
					},
				},
				nil,
			},
			{
				"3(B)3(1)/LDP",
				models.EventDetail{
					Play:       models.LinedIntoDoublePlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionFirstBase,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierLinedIntoDoublePlay,
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},

			{
				"64(1)3/GDP/G6",
				models.EventDetail{
					Play:       models.GroundedIntoDoublePlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionShortStop,
						models.PositionSecondBase,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBallDoublePlay,
						},
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "6",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"4(1)3/G4/GDP",
				models.EventDetail{
					Play:       models.GroundedIntoDoublePlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionSecondBase,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "4",
						},
						{
							PlayModifier: models.ModifierGroundBallDoublePlay,
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"1(B)16(2)63(1)/LTP/L1",
				models.EventDetail{
					Play:       models.LinedIntoTriplePlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionPitcher,
						models.PositionPitcher,
						models.PositionShortStop,
						models.PositionShortStop,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierLinedIntoTriplePlay,
						},
						{
							PlayModifier: models.ModifierLinedDrive,
							Location:     "1",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"D7/G5.3-H;2-H;1-H",
				models.EventDetail{
					Play:       models.Double,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionLeftField,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "5",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
						{
							StartBase:  2,
							FinishBase: 4,
						},
						{
							StartBase:  1,
							FinishBase: 4,
						},
					},
				},
				nil,
			},
			{
				"HR/F78XD.2-H;1-H",
				models.EventDetail{
					Play:       models.HomeRun,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierFlyBall,
							Location:     "78XD",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  2,
							FinishBase: 4,
						},
						{
							StartBase:  1,
							FinishBase: 4,
						},
					},
				},
				nil,
			},
			{
				"HR9/F9LS.3-H;1-H",
				models.EventDetail{
					Play:       models.HomeRun,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierFlyBall,
							Location:     "9LS",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
						{
							StartBase:  1,
							FinishBase: 4,
						},
					},
				},
				nil,
			},
			{
				"FC5/G5.3XH(52)",
				models.EventDetail{
					Play:       models.FieldersChoice,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionThirdBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "5",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"8/F78",
				models.EventDetail{
					Play:       models.FlyBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionCenterField,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierFlyBall,
							Location:     "78",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},

			{
				"FLE5/P5F",
				models.EventDetail{
					Play:       models.ErrorOnFlyBall,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionThirdBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierPopup,
							Location:     "5F",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},

			{
				"E3.1-2;B-1",
				models.EventDetail{
					Play:       models.Error,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
						{
							StartBase:  0,
							FinishBase: 1,
						},
					},
				},
				nil,
			},

			{
				"K23+WP.2-3",
				models.EventDetail{
					Play: models.StrikeOut,
					ExtraPlays: []models.BasicPlay{
						models.WildPitch,
					},
					Fielders: []models.Position{
						models.PositionCatcher,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  2,
							FinishBase: 3,
						},
					},
				},
				nil,
			},
			{
				"W+WP.2-3",
				models.EventDetail{
					Play: models.Walk,
					ExtraPlays: []models.BasicPlay{
						models.WildPitch,
					},
					Fielders:  []models.Position{},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  2,
							FinishBase: 3,
						},
					},
				},
				nil,
			},
			{
				"CSH(12)",
				models.EventDetail{
					Play:       models.CaughtStealing,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionPitcher,
						models.PositionCatcher,
					},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"BK.3-H;1-2",
				models.EventDetail{
					Play:       models.Balk,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},

			{
				"IW",
				models.EventDetail{
					Play:       models.IntentionalWalk,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv:  []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"W.1-2",
				models.EventDetail{
					Play:       models.Walk,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},
			{
				"NP",
				models.EventDetail{
					Play:       models.NoPlay,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv:  []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"K+WP.B-1",
				models.EventDetail{
					Play: models.StrikeOut,
					ExtraPlays: []models.BasicPlay{
						models.WildPitch,
					},
					Fielders:  []models.Position{},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  0,
							FinishBase: 1,
						},
					},
				},
				nil,
			},
			{
				"K+PB.1-2",
				models.EventDetail{
					Play: models.StrikeOut,
					ExtraPlays: []models.BasicPlay{
						models.PassedBall,
					},
					Fielders:  []models.Position{},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},
			{
				"K23",
				models.EventDetail{
					Play:       models.StrikeOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionCatcher,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"K",
				models.EventDetail{
					Play:       models.StrikeOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv:  []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"HP.1-2",
				models.EventDetail{
					Play:       models.HitByPitch,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers:  []models.Modifier{},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},

			{
				"H/L7D",
				models.EventDetail{
					Play:       models.HomeRun,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierLinedDrive,
							Location:     "7D",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},

			{
				"FC3/G3S.3-H;1-2",
				models.EventDetail{
					Play:       models.FieldersChoice,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "3S",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},

			{
				"E1/TH/BG15.1-3",
				models.EventDetail{
					Play:       models.Error,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionPitcher,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierThrowing,
						},
						{
							PlayModifier: models.ModifierGroundBallBunt,
							Location:     "15",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 3,
						},
					},
				},
				nil,
			},
			{
				"T9/F9LD.2-H",
				models.EventDetail{
					Play:       models.Triple,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionRightField,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierFlyBall,
							Location:     "9LD",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  2,
							FinishBase: 4,
						},
					},
				},
				nil,
			},

			{
				"S9/G",
				models.EventDetail{
					Play:       models.Single,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionRightField,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"C/E1.1-2",
				models.EventDetail{
					Play:       models.CatcherInterference,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierErrorOn,
							Location:     "1",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},
			{
				"C/E2.1-2",
				models.EventDetail{
					Play:       models.CatcherInterference,
					ExtraPlays: []models.BasicPlay{},
					Fielders:   []models.Position{},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierErrorOn,
							Location:     "2",
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},

			{
				"23/SH.1-2",
				models.EventDetail{
					Play:       models.GroundBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionCatcher,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierSacrificeBunt,
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},

			{
				"54(B)/BG25/SH.1-2",
				models.EventDetail{
					Play:       models.GroundBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionThirdBase,
						models.PositionSecondBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBallBunt,
							Location:     "25",
						},
						{
							PlayModifier: models.ModifierSacrificeBunt,
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  1,
							FinishBase: 2,
						},
					},
				},
				nil,
			},
			{
				"143/G1",
				models.EventDetail{
					Play:       models.GroundBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionPitcher,
						models.PositionSecondBase,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "1",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"9/SF.3-H",
				models.EventDetail{
					Play:       models.FlyBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionRightField,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierSacrificeFly,
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  3,
							FinishBase: 4,
						},
					},
				},
				nil,
			},
			{
				"63/G6M",
				models.EventDetail{
					Play:       models.GroundBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionShortStop,
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
							Location:     "6M",
						},
					},
					RunnerAdv: []models.RunnerAdvance{},
				},
				nil,
			},
			{
				"3/G.2-3",
				models.EventDetail{
					Play:       models.FlyBallOut,
					ExtraPlays: []models.BasicPlay{},
					Fielders: []models.Position{
						models.PositionFirstBase,
					},
					Modifiers: []models.Modifier{
						{
							PlayModifier: models.ModifierGroundBall,
						},
					},
					RunnerAdv: []models.RunnerAdvance{
						{
							StartBase:  2,
							FinishBase: 3,
						},
					},
				},
				nil,
			},
		}
		for _, test := range tests {
			convey.Convey("Parse "+test.t+"...", func() {
				ed, ok := readers.ParseEventDetail(test.t)
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(ed, convey.ShouldResemble, test.v)
			})
		}
	})
}
