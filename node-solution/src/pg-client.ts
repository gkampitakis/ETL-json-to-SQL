import pg_promise from 'pg-promise';
import { pg_config } from './configurator';
import { Matchup } from './types';

const pgp = pg_promise({
  capSQL: true,
});
const db = pgp(pg_config.client);
const cs = new pgp.helpers.ColumnSet([
  'gold_earned',
  'minions_killed',
  'kda',
  'champion',
  'vision_score',
  'summoner_name',
  'win',
  'game_version',
  'damage_dealt_to_champions',
  'lane',
  'region'
], { table: pg_config.table });

export function bulkInsert(data: Matchup[]) {
  const query = () => pgp.helpers.insert(data, cs);
  return db.none(query, data);
}
