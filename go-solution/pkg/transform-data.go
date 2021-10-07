package pkg

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

// FIXME: add tags
type Matchup struct {
	GoldEarned             int
	MinionsKilled          int
	KDA                    int
	Champion               string
	VisionScore            int
	SummoreName            string
	Win                    bool
	GameVersion            string
	DamageDealtToChampions int
	Lane                   string
	Region                 string
}

// Returns a mathcup and isValid value
// If isValid is false matchup is going to be empty
func TransformDatum(datum string) (Matchup, bool) {
	commaRegex := regexp.MustCompile(`^,*`)
	if datum == "[" || datum == "]" {
		return Matchup{}, false
	}

	var parsedDatum map[string]interface{}
	err := json.Unmarshal(commaRegex.ReplaceAll([]byte(datum), []byte("")), &parsedDatum)
	if err != nil {
		log.Println(err)
		return Matchup{}, false
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

	return Matchup{
		Champion:               parsedDatum["champion"].(string),
		DamageDealtToChampions: int(parsedDatum["totaldamagedealttochampions"].(float64)),
		GoldEarned:             int(parsedDatum["goldearned"].(float64)),
		Win:                    parsedDatum["win"] == "true",
		MinionsKilled:          int(parsedDatum["totalminionskilled"].(float64)),
		KDA:                    getKda(kills, assists, deaths),
		Lane:                   lane,
		Region:                 region,
		SummoreName:            parsedDatum["summonername"].(string),
		VisionScore:            int(parsedDatum["visionscore"].(float64)),
	}, true
}

func getKda(kills, assists, deaths int) int {
	if deaths == 0 {
		deaths = 1
	}

	return (kills + assists) / deaths
}
