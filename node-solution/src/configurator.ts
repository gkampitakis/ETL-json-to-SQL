import dotenv from 'dotenv';

if (process.env.NODE_ENV !== 'production') {
  dotenv.config();
}

export const config = {
  filePath: process.env.FILE_PATH ?? '../data-to-load/matchups.json',
  batchRecords: parseInt(process.env.BATCH_RECORDS ?? '5000'),
  errorLogPath: process.env.ERROR_LOG_PATH ?? `${process.cwd()}/errors`
};

export const pg_config = {
  client: {
    keepAlive: true,
    allowExitOnIdle: true,
    host: process.env.PG_HOST ?? 'localhost',
    user: process.env.PG_USER ?? 'ETL_user',
    password: process.env.PG_PASS ?? 'ETL_pass',
    database: process.env.PG_DATABASE ?? 'ETL_db',
    port: parseInt(process.env.PG_PORT ?? '5432'),
  },
  table: process.env.PG_TABLE ?? 'matchups'
}
