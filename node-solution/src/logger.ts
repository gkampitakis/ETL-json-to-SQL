import Logger from 'pino';

export default Logger({
  level: process.env.LOG_LEVEL ?? 'debug',
  name: 'ETL'
});