package models

import (
	"encoding/json"
	"fmt"

	"github.com/gocraft/dbr"
)

type EventType int

const (
	GameID     EventType = 0
	Version    EventType = 1
	Info       EventType = 2
	Start      EventType = 3
	Play       EventType = 4
	Sub        EventType = 5
	Data       EventType = 6
	Comment    EventType = 7
	BatterAdj  EventType = 8
	LineupAdj  EventType = 9
	PitcherAdj EventType = 10
)

type InfoType int

const (
	InputProgVers InfoType = 0
	VisitingTeam  InfoType = 1
	HomeTeam      InfoType = 2
	GameDate      InfoType = 3

	// Add other Info types if needed
)

type InningHalf int

const (
	TopHalf    InningHalf = 0
	BottomHalf InningHalf = 1
)

type BasicPlay int

const (
	FlyBallOut BasicPlay = iota
	GroundBallOut
	GroundedIntoDoublePlay
	LinedIntoDoublePlay
	LinedIntoTriplePlay
	CatcherInterference
	Single
	Double
	Triple
	HomeRun
	GroundRuleDouble
	Error
	FieldersChoice
	ErrorOnFlyBall
	HitByPitch
	StrikeOut
	NoPlay
	Walk
	IntentionalWalk
	Balk
	CaughtStealing
	DefensiveIndifference
	OtherAdvance
	PassedBall
	WildPitch
	PickOff
	PickOffCaughtStealing
	StolenBase
)

func (bp BasicPlay) String() string {
	names := [...]string{
		"FlyBallOut",
		"GroundBallOut",
		"GroundedIntoDoublePlay",
		"LinedIntoDoublePlay",
		"LinedIntoTriplePlay",
		"CatcherInterference",
		"Single",
		"Double",
		"Triple",
		"HomeRun",
		"GroundRuleDouble",
		"Error",
		"FieldersChoice",
		"ErrorOnFlyBall",
		"HitByPitch",
		"StrikeOut",
		"NoPlay",
		"Walk",
		"IntentionalWalk",
		"Balk",
		"CaughtStealing",
		"DefensiveIndifference",
		"OtherAdvance",
		"PassedBall",
		"WildPitch",
		"PickOff",
		"PickOffCaughtStealing",
		"StolenBase",
	}

	if bp < FlyBallOut || bp > StolenBase {
		return "Invalid basic play"
	}
	return names[bp]
}

type PlayModifier int

const (
	ModifierFlyBall PlayModifier = iota
	ModifierLinedDrive
	ModifierErrorOn

	ModifierThrowing

	ModifierSacrificeFly
	ModifierSacrificeBunt

	ModifierUnassisted

	ModifierGroundBall
	ModifierGroundBallBunt
	ModifierGroundBallDoublePlay

	ModifierLinedIntoDoublePlay
	ModifierLinedIntoTriplePlay

	ModifierPopup

	ModifierGroundBallDoublePlayBunt
	ModifierUmpireReviewCallOnField
	ModifierUmpireInterference
	ModifierUnspecifiedTriplePlay
	ModifierRelayThrow
	ModifierPassedRunner
	ModifierObstruction
	ModifierNoDoublePlay
	ModifierManagerChallenge
	ModifierInsideTheParkHomeRun
	ModifierInterference
	ModifierInfieldFlyRule
	ModifierGroundBallTriplePlay
	ModifierFoulBall
	ModifierFanInterference
	ModifierFlyBallDoublePlay
	ModifierUnspecifiedDoublePlay
	ModifierCourtesyBatter
	ModifierCourtesyFielder
	ModifierCourtesyRunner
	ModifierCalledThirdStrike
	ModifierRunnerHitByBattedBall
	ModifierPopupDoublePlayBunt
	ModifierBattingOutOfTurn
	ModifierLinedDriveBunt
	ModifierBatterInterference
	ModifierAppealPlay
	ModifierPopupBunt
	ModifierForceOut
)

func (pm PlayModifier) String() string {
	names := [...]string{
		"Fly ball",
		"Line drive",
		"Error",
		"Throwing",
		"Sacrifice fly",
		"Sacrifice bunt",
		"Unassisted",
		"Ground ball",
		"Ground ball bunt",
		"Ground ball double play",
		"Lined into double play",
		"Lined into triple play",
		"Popup",
		"Ground ball double play bunt",
		"Umpire review call on field",
		"Umpire interference",
		"Unspecified triple play",
		"Relay throw",
		"Runner passed another runner",
		"Obstruction",
		"No double play",
		"Manager challenge",
		"Inside the park homerun",
		"Interference",
		"Infield fly rule",
		"Ground ball triple play",
		"Foul ball",
		"Fan interference",
		"Fly ball double play",
		"Unspecified double play",
		"Courtesy batter",
		"Courtesy fielder",
		"Courtesy runner",
		"Called third strike",
		"Runner hit by batted ball",
		"Popup double play bunt",
		"Batting out of turn",
		"Lined drive bunt",
		"Batter interference",
		"Appeal play",
		"Popup bunt",
		"Forced out",
	}

	if pm < ModifierFlyBall || pm > ModifierForceOut {
		return "Invalid play modifier"
	}
	return names[pm]
}

type HitLocation string

type Modifier struct {
	PlayModifier
	Location HitLocation
}

func (m Modifier) String() string {
	return fmt.Sprintf("%s, %s", m.PlayModifier, m.Location)
}

type RunnerAdvance struct {
	StartBase  int
	FinishBase int
}

func (ra RunnerAdvance) String() string {
	return fmt.Sprintf("%d to %d", ra.StartBase, ra.FinishBase)
}

func NewGameEvent(gameID int, et EventType, inning int, half InningHalf, player int) GameEvent {
	ge := GameEvent{
		GameID:     gameID,
		Event:      et,
		Inning:     inning,
		InningHalf: half,
		Player:     player,
	}
	return ge
}

func (ed EventDetail) String() string {
	return fmt.Sprintf("\nPlay: %s\n\tFielders: %s\n\tModifiers: %s\n\tRunners: %s\n",
		ed.Play, ed.Fielders, ed.Modifiers, ed.RunnerAdv)
}

type GameEvent struct {
	ID         int         `db:"id" json:"-"`
	GameID     int         `db:"game_id" json:"-"`
	Event      EventType   `db:"event" json:"-"`
	Inning     int         `db:"inning" json:"-"`
	InningHalf InningHalf  `db:"inning_half" json:"-"`
	Player     int         `db:"player_id" json:"-"`
	Play       EventDetail `db:"-" json:"event_detail"`
	PlayJSON   string      `db:"event_detail" json:"-"`
}

type EventDetail struct {
	Play       BasicPlay
	ExtraPlays []BasicPlay
	Fielders   []Position
	Modifiers  []Modifier
	RunnerAdv  []RunnerAdvance
}

func (g GameEvent) String() string {
	return fmt.Sprintf("%d, %d (%d, %d)\n\tPlayer: %d\n\t%s", g.GameID, g.Event, g.Inning, g.InningHalf, g.Player, g.Play)
}

func (g *GameEvent) marshalEventDetail() error {
	b, err := json.Marshal(g)
	if err != nil {
		return err
	}
	g.PlayJSON = string(b)
	return nil
}

func (e *GameEvent) Save(session dbr.SessionRunner) error {
	e.marshalEventDetail()
	_, err := session.InsertInto("game_events").
		Columns("game_id", "player_id", "event", "inning", "inning_half", "event_detail").
		Record(e).
		Exec()
	return err
}

func SaveGamesEvents(session dbr.SessionRunner, games []GameEvent) error {
	var err error
	for _, g := range games {
		err = g.Save(session)
		if err != nil {
			break
		}
	}
	return err
}
