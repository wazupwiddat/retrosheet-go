package readers

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wazupwiddat/retrosheet/models"
)

var (
	eventPatternMap = map[string]models.BasicPlay{
		`^[0-9]{1}`:                                models.FlyBallOut,
		`^[0-9]{2,}(\((.)\))?`:                     models.GroundBallOut,
		`^[0-9]+\((.)\)[0-9]{1}`:                   models.GroundedIntoDoublePlay,
		`^[0-9]+\((.)\)[0-9]+\((.)\)`:              models.LinedIntoDoublePlay,
		`^[0-9]+\((.)\)[0-9]+\((.)\)[0-9]+\((.)\)`: models.LinedIntoTriplePlay,
		`^C`:           models.CatcherInterference,
		`^S[0-9]`:      models.Single,
		`^D[0-9]`:      models.Double,
		`^T[0-9]`:      models.Triple,
		`^HP`:          models.HitByPitch,
		`^H[R]?[0-9]?`: models.HomeRun,
		`^DGR?`:        models.GroundRuleDouble,
		`^E([0-9])?`:   models.Error,
		`^FC([0-9])?`:  models.FieldersChoice,
		`^FLE([0-9])?`: models.ErrorOnFlyBall,
		`^K([0-9]+)?(\+WP|\+PB)?(\+SB[0-9])?(\+CS[0-9])?(\+PO[0-9])?(\+E[0-9])?`: models.StrikeOut,
		`^NP`: models.NoPlay,
		`^W(\+WP|\+PB)?(\+SB[0-9])?(\+CS[0-9])?(\+PO[0-9])?(\+E[0-9])?`:  models.Walk,
		`^IW(\+WP|\+PB)?(\+SB[0-9])?(\+CS[0-9])?(\+PO[0-9])?(\+E[0-9])?`: models.IntentionalWalk,
		`^BK`: models.Balk,
		`^CS[2-3H]?\(([0-9A-Z]+)\)?`: models.CaughtStealing,
		`^DI`: models.DefensiveIndifference,
		`^OA`: models.OtherAdvance,
		`^PB`: models.PassedBall,
		`^WP`: models.WildPitch,
		`^PO[1-3]\([0-9A-Z]+\)`:        models.PickOff,
		`^POCS[2-3H]?\(([0-9A-Z]+)\)?`: models.PickOffCaughtStealing,
		`^(SB[2-3H][;]?)+`:             models.StolenBase,
	}
	parseEventTypeMap = map[string]models.EventType{
		"id":      models.GameID,
		"version": models.Version,
		"info":    models.Info,
		"start":   models.Start,
		"play":    models.Play,
		"sub":     models.Sub,
		"data":    models.Data,
		"com":     models.Comment,
		"badj":    models.BatterAdj,
		"ladj":    models.LineupAdj,
		"padj":    models.PitcherAdj,
	}
	parseInfoTypeMap = map[string]models.InfoType{
		"visteam":  models.VisitingTeam,
		"hometeam": models.HomeTeam,
		"date":     models.GameDate,
	}
	parseInningHalfMap = map[string]models.InningHalf{
		"0": models.TopHalf,
		"1": models.BottomHalf,
	}
	positionMap = map[string]models.Position{
		"1": models.PositionPitcher,
		"2": models.PositionCatcher,
		"C": models.PositionCatcher,
		"3": models.PositionFirstBase,
		"4": models.PositionSecondBase,
		"5": models.PositionThirdBase,
		"6": models.PositionShortStop,
		"7": models.PositionLeftField,
		"8": models.PositionCenterField,
		"9": models.PositionRightField,
	}
	playModifierMap = map[string]models.PlayModifier{
		"^AP$": models.ModifierAppealPlay,
		"^BP([0-9]|$)([0-9A-Z]+)?": models.ModifierPopupBunt,
		"^BG([0-9]|$)([0-9A-Z]+)?": models.ModifierGroundBallBunt,
		"^BGDP$":                   models.ModifierGroundBallDoublePlayBunt,
		"^BINT":                    models.ModifierBatterInterference,
		"^BL([0-9]|$)([0-9A-Z]+)?": models.ModifierLinedDriveBunt,
		"^BOOT$":                   models.ModifierBattingOutOfTurn,
		"^BPDP$":                   models.ModifierPopupDoublePlayBunt,
		"^BR$":                     models.ModifierRunnerHitByBattedBall,
		"^C$":                      models.ModifierCalledThirdStrike,
		"^COUB$":                   models.ModifierCourtesyBatter,
		"^COUF$":                   models.ModifierCourtesyFielder,
		"^COUR$":                   models.ModifierCourtesyRunner,
		"^DP$":                     models.ModifierUnspecifiedDoublePlay,
		"^F([0-9]|$)([0-9A-Z]+)?": models.ModifierFlyBall,
		"^E([0-9])?":              models.ModifierErrorOn,
		"^FDP$":                   models.ModifierFlyBallDoublePlay,
		"^FINT$":                  models.ModifierFanInterference,
		"^FL$":                    models.ModifierFoulBall,
		"^FO$":                    models.ModifierForceOut,
		"^G([0-9]|$)([0-9A-Z]+)?": models.ModifierGroundBall,
		"^GDP$":                   models.ModifierGroundBallDoublePlay,
		"^GTP$":                   models.ModifierGroundBallTriplePlay,
		"^IF$":                    models.ModifierInfieldFlyRule,
		"^INT$":                   models.ModifierInterference,
		"^IPHR$":                  models.ModifierInsideTheParkHomeRun,
		"^L([0-9]|$)([0-9A-Z]+)?": models.ModifierLinedDrive,
		"^LDP$":                   models.ModifierLinedIntoDoublePlay,
		"LTP$":                    models.ModifierLinedIntoTriplePlay,
		"^MREV$":                  models.ModifierManagerChallenge,
		"^NDP$":                   models.ModifierNoDoublePlay,
		"^OBS$":                   models.ModifierObstruction,
		"^P([0-9]|$)([0-9A-Z]+)?": models.ModifierPopup,
		"^PASS$":                  models.ModifierPassedRunner,
		"^R([0-9])?":              models.ModifierRelayThrow,
		"^SF$":                    models.ModifierSacrificeFly,
		"^SH$":                    models.ModifierSacrificeBunt,
		"^TH([0-9])?":             models.ModifierThrowing,
		"^TP$":                    models.ModifierUnspecifiedTriplePlay,
		"^UINT$":                  models.ModifierUmpireInterference,
		"^UREV$":                  models.ModifierUmpireReviewCallOnField,
	}
	runnerAdvanceMap = map[string]models.RunnerAdvance{
		"^B-1[#]?": {StartBase: 0, FinishBase: 1},
		"^B-2[#]?": {StartBase: 0, FinishBase: 2},
		"^B-3[#]?": {StartBase: 0, FinishBase: 3},
		"^B-H[#]?": {StartBase: 0, FinishBase: 4},
		"^1-2[#]?": {StartBase: 1, FinishBase: 2},
		"^1-3[#]?": {StartBase: 1, FinishBase: 3},
		"^1-H[#]?": {StartBase: 1, FinishBase: 4},
		"^2-3[#]?": {StartBase: 2, FinishBase: 3},
		"^2-H[#]?": {StartBase: 2, FinishBase: 4},
		"^3-H[#]?": {StartBase: 3, FinishBase: 4},
	}
)

func ParseEventType(val string) (models.EventType, bool) {
	i, ok := parseEventTypeMap[val]
	if !ok {
		return -1, ok
	}
	return i, ok
}

func ParseInfoType(val string) (models.InfoType, bool) {
	i, ok := parseInfoTypeMap[val]
	if !ok {
		return -1, ok
	}
	return i, ok
}

func ParseLeague(val string) models.League {
	if val == "A" || val == "a" {
		return models.American
	}
	if val == "N" || val == "n" {
		return models.National
	}
	return -1
}

func ParseHanded(val string) models.Handed {
	if val == "R" || val == "r" {
		return models.RightHanded
	}
	if val == "L" || val == "l" {
		return models.LeftHanded
	}
	if val == "B" || val == "b" {
		return models.BothHanded
	}
	return -1
}

func ParseYear(val string) int {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(val, "")
	i, err := strconv.Atoi(processedString)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ParseInning(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}
	return i
}

func ParseInningHalf(val string) (models.InningHalf, bool) {
	i, ok := parseInningHalfMap[val]
	if !ok {
		return -1, ok
	}
	return i, ok

}

func ParseBasicPlay(val string) models.BasicPlay {
	for regx, result := range eventPatternMap {
		reg, err := regexp.Compile(regx)
		if err != nil {
			log.Fatal(err)
		}
		if v := reg.FindString(val); v == val {
			return result
		}
	}

	return models.NoPlay
}

func parseFieldersInPlay(val string, play models.BasicPlay) []models.Position {
	var fieldersStr string
	switch play {
	case models.LinedIntoTriplePlay, models.LinedIntoDoublePlay, models.GroundedIntoDoublePlay:
		fieldersStr = removeStringRegex(`(\([0-9]+\))|(\(B)\)`, val)
	case models.CaughtStealing:
		fieldersStr = matchStringRegex(`(\([0-9]+\))`, val)
		fieldersStr = matchStringRegex(`[0-9]+`, val)
	case models.HomeRun:
		// do nothing here
	default:
		fieldersStr = matchStringRegex(`[0-9]+`, val)
	}

	fielders := []models.Position{}
	for _, f := range fieldersStr {
		i, _ := strconv.Atoi(string(f))
		fielders = append(fielders, models.Position(i))
	}
	return fielders
}

func parseExtraEvents(val string, play models.BasicPlay) []models.BasicPlay {
	extras := []models.BasicPlay{}
	switch play {
	case models.Walk, models.IntentionalWalk, models.StrikeOut:
		mods := strings.Split(val, "+")
		if len(mods) > 1 {
			for _, mod := range mods[1:] {
				extras = append(extras, ParseBasicPlay(mod))
			}
		}
	}
	return extras
}

func parsePlayMod(vals []string, play models.BasicPlay) []models.Modifier {
	mods := []models.Modifier{}
	for _, val := range vals {
		for regx, result := range playModifierMap {
			reg, err := regexp.Compile(regx)
			if err != nil {
				log.Fatal(err)
			}
			if v := reg.FindString(val); v == val {
				mod := models.Modifier{
					PlayModifier: result,
				}

				mod.Location = parseModifierHitLocation(val, result)
				mods = append(mods, mod)
				break
			}
		}
	}

	return mods
}

func parseModifierHitLocation(val string, mod models.PlayModifier) models.HitLocation {
	loc := matchStringRegex(`[0-9](.+)?`, val)
	return models.HitLocation(loc)
}

func ParseRunnerAdvances(vals []string) []models.RunnerAdvance {
	runners := []models.RunnerAdvance{}
	for _, r := range vals {
		just := removeStringRegex(`(\(.+\))`, r)
		for regx, result := range runnerAdvanceMap {
			reg, err := regexp.Compile(regx)
			if err != nil {
				log.Fatal(err)
			}
			if v := reg.FindString(just); v == just {
				runners = append(runners, result)
				break
			}
		}
		if len(runners) == 0 {
			// log.Println("Unknown runner advance: ", r)
		}

	}
	return runners
}

func ParseEventDetail(val string) (models.EventDetail, bool) {
	// <basicplay>/<basicplaymodifier>.<runners>

	eventDetail := models.EventDetail{}

	// <basicplay>/<basicplaymodifier>
	basicPlay, modifiers, ok := splitPlayWithModifier(val)
	eventDetail.Play = ParseBasicPlay(basicPlay)
	eventDetail.Fielders = parseFieldersInPlay(basicPlay, eventDetail.Play)

	// events on walks and strikeouts
	eventDetail.ExtraPlays = parseExtraEvents(basicPlay, eventDetail.Play)

	eventDetail.Modifiers = parsePlayMod(modifiers, eventDetail.Play)

	// multiple runners are separated by ;
	// runners advance safely are indicated with a - (2-3)
	// runners thrown out are indicated with a X (2X3)

	runners, ok := splitRunners(val)
	if !ok {
		log.Println("Unknown runner scenario: ", basicPlay)
		return eventDetail, ok
	}

	eventDetail.RunnerAdv = ParseRunnerAdvances(runners)

	return eventDetail, ok
}

func splitPlayWithModifier(val string) (string, []string, bool) {
	playWithModifier := strings.Split(val, ".")
	if len(playWithModifier) > 0 {
		s := strings.Split(playWithModifier[0], "/")
		return s[0], s[1:], true
	}
	return "NO PLAY", []string{}, false
}

func splitRunners(val string) ([]string, bool) {
	runners := strings.Split(val, ".")
	if len(runners) > 1 {
		s := strings.Split(runners[1], ";")
		return s[0:], true
	}
	return []string{}, true
}

func removeStringRegex(regex string, val string) string {
	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(val, "")
}
func matchStringRegex(regex string, val string) string {
	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}
	return reg.FindString(val)
}
