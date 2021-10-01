import { createStreamFileReader } from './src/file-reader';
import { bulkInsert } from './src/pg-client';
import { batchProcessing } from './src/batch-processor';
import Logger from './src/logger';

async function main() {
  Logger.info('ETL pipeline starting 🚀');

  const {
    getData,
    pause,
    resume,
  } = createStreamFileReader('../data-to-load/matchups.json');

  batchProcessing({
    getData,
    pause,
    resume,
    commitData: bulkInsert
  });

  Logger.info('ETL pipeline finished 🤖');
}

main();
