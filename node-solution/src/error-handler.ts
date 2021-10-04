import { appendFile, mkdir } from 'fs/promises';
import { json2csvAsync } from 'json-2-csv';
import { Matchup } from './types';
import Logger from './logger';

export async function saveToErrorLog(data: Matchup[], path: string) {
  const serializedData = await json2csvAsync(data, {
    checkSchemaDifferences: true,
    prependHeader: false
  });

  await mkdir(path, {
    recursive: true,
  });

  return appendFile(`${path}/log.csv`, serializedData + '\n', {
    flag: 'a+',
    encoding: 'utf-8'
  })
    .catch((error) => {
      Logger.error(`[Save Logger]:${error}`);
    });
}