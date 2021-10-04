import { createStreamFileReader } from './src/file-reader';
import { bulkInsert } from './src/pg-client';
import { batchProcessing } from './src/batch-processor';
import { config } from './src/configurator';
import Logger from './src/logger';

const { filePath, ...batchProcessingConfig } = config;

async function main(filePath: string) {
  Logger.info('ETL pipeline starting 🚀');
  console.log('[Starting with config]: ', config);

  const {
    getData,
    pause,
    resume,
    onEnd
  } = createStreamFileReader(filePath);

  const batchProcessor = batchProcessing({
    getData,
    pause,
    resume,
    onEnd,
    bulkInsert
  }, batchProcessingConfig);

  batchProcessor.on('finish', (report) => {
    console.log('[Report]: ', report);
    Logger.info('ETL pipeline finished 🤖');
  });
}

main(filePath);
