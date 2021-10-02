import { appendFile, mkdir } from 'fs/promises';
import { json2csvAsync } from 'json-2-csv';
import { Matchup } from './types';
import Logger from './logger';

const ERROR_LOG_PATH = process.env.ERROR_LOG_PATH ?? `${process.cwd()}/errors`;

export async function saveToErrorLog(data: Matchup[]) {
  const serializedData = await json2csvAsync(data, {
    checkSchemaDifferences: true,
    prependHeader: false,
  });

  await mkdir(ERROR_LOG_PATH, {
    recursive: true,
  });

  return appendFile(`${ERROR_LOG_PATH}/log.csv`, serializedData, {
    flag: 'a+',
    encoding: 'utf-8'
  })
    .catch((error) => {
      Logger.error(`[Save Logger]:${error}`);
    });
}