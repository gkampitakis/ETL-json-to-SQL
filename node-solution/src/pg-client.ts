import pg_promise from 'pg-promise'
import { Matchup } from './types';
const pgp = pg_promise({
  capSQL: true,
});

const db = pgp({
  keepAlive: true,
  host: process.env.PG_HOST ?? 'localhost',
  user: process.env.PG_USER ?? 'ETL_user',
  password: process.env.PG_PASS ?? 'ETL_pass',
  database: process.env.PG_DATABASE ?? 'ETL_db',
  port: parseInt(process.env.PG_PORT ?? '5432'),
  allowExitOnIdle: true
});

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
], { table: process.env.PG_TABLE ?? 'matchups' });

export function bulkInsert(data: Matchup[]) {
  const query = () => pgp.helpers.insert(data, cs);
  return db.none(query, data);
}
