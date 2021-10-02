import { createStreamFileReader } from './src/file-reader';
import { bulkInsert } from './src/pg-client';
import { batchProcessing } from './src/batch-processor';
import Logger from './src/logger';

async function main(filePath: string) {
  Logger.info('ETL pipeline starting 🚀');

  const {
    getData,
    pause,
    resume,
    onEnd
  } = createStreamFileReader(filePath);

  const { on } = batchProcessing({
    getData,
    pause,
    resume,
    onEnd,
    bulkInsert
  });

  on('finish', (report) => {
    Logger.info('ETL pipeline finished 🤖');
  });
}

main(process.env.FILE_PATH ?? '../data-to-load/matchups.json');
