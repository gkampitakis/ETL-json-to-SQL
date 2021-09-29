import { Pool } from 'pg';

export const client = new Pool({
  keepAlive: true,
  host: process.env.PG_HOST ?? 'localhost',
  user: process.env.PG_PASS ?? 'ETL_user',
  password: process.env.PG_PASS ?? 'ETL_pass',
  database: process.env.PG_DATABASE ?? 'ETL_db',
  port: parseInt(process.env.PG_PORT ?? '5432')
});
