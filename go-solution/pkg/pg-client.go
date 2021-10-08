package pkg

import (
	"context"

	"github.com/jackc/pgx/v4"
)

var columns = []string{
	"gold_earned",
	"minions_killed",
	"kda",
	"champion",
	"vision_score",
	"summoner_name",
	"win",
	"game_version",
	"damage_dealt_to_champions",
	"lane",
	"region",
}

func BulkInsert(conn *pgx.Conn, table string, data []Matchup) (int64, error) {
	return conn.CopyFrom(
		context.Background(),
		pgx.Identifier{table},
		columns,
		pgx.CopyFromSlice(len(data), func(idx int) ([]interface{}, error) {
			return []interface{}{
					data[idx].GoldEarned,
					data[idx].MinionsKilled,
					data[idx].KDA,
					data[idx].Champion,
					data[idx].VisionScore,
					data[idx].SummonerName,
					data[idx].Win,
					data[idx].GameVersion,
					data[idx].DamageDealtToChampions,
					data[idx].Lane,
					data[idx].Region,
				},
				nil
		}))
}
