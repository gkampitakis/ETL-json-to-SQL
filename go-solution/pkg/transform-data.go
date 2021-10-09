package pkg

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

type Matchup struct {
	GoldEarned             int
	MinionsKilled          int
	KDA                    int
	Champion               string
	VisionScore            int
	SummonerName           string
	Win                    bool
	GameVersion            string
	DamageDealtToChampions int
	Lane                   string
	Region                 string
}

// Returns a mathcup, erred and and ignore.
// If ignore is true or erred is true, matchup is going to be empty
func TransformDatum(datum string) (m Matchup, erred bool, ignore bool) {
	// NOTE: to future self by having named return variables, defer function can manipulate those values
	defer func() {
		if e := recover(); e != nil {
			m = Matchup{}
			erred = true
			ignore = true
			return
		}
	}()
	commaRegex := regexp.MustCompile(`^,*`)
	if datum == "[" || datum == "]" {
		return Matchup{}, false, true
	}

	var parsedDatum map[string]interface{}
	err := json.Unmarshal(commaRegex.ReplaceAll([]byte(datum), []byte("")), &parsedDatum)
	if err != nil {
		log.Println(err)
		return Matchup{}, true, true
	}

	kills := int(parsedDatum["kills"].(float64))
	assists := int(parsedDatum["assists"].(float64))
	deaths := int(parsedDatum["deaths"].(float64))
	region, lane := func() (string, string) {
		tmp := strings.Split(parsedDatum["p_match_id"].(string), "_")

		if tmp[2] == "utility" {
			return tmp[0], "support"
		}

		return tmp[0], tmp[2]
	}()
	// NOTE: to future self, when type asserting, the second value
	// is a boolean indicating if the assertion was correct or not, saving from unhandled panics
	summonerName, _ := parsedDatum["summonername"].(string)

	return Matchup{
		Champion:               parsedDatum["champion"].(string),
		DamageDealtToChampions: int(parsedDatum["totaldamagedealttochampions"].(float64)),
		GoldEarned:             int(parsedDatum["goldearned"].(float64)),
		Win:                    parsedDatum["win"] == "true",
		GameVersion:            parsedDatum["gameversion"].(string),
		MinionsKilled:          int(parsedDatum["totalminionskilled"].(float64)),
		KDA:                    getKda(kills, assists, deaths),
		Lane:                   lane,
		Region:                 region,
		SummonerName:           summonerName,
		VisionScore:            int(parsedDatum["visionscore"].(float64)),
	}, false, false
}

func getKda(kills, assists, deaths int) int {
	if deaths == 0 {
		deaths = 1
	}

	return (kills + assists) / deaths
}
