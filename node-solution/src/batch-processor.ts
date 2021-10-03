import {
  Matchup,
  BatchProcessingParams,
  BatchProcessingEventCallback
} from './types';
import { saveToErrorLog } from './error-handler';
import Logger from './logger';

const BATCH_RECORDS = parseInt(process.env.BATCH_RECORDS ?? '5000');

export function batchProcessing({
  getData,
  pause,
  bulkInsert,
  resume,
  onEnd
}: BatchProcessingParams): { on: BatchProcessingEventCallback } {
  const events: Record<string, any> = {};
  const promises: Promise<void>[] = [];
  const startTime = Date.now();
  let buffer: Matchup[] = [];
  let rowsInserted = 0;
  let bulkInsertErrors = 0;
  let transformationErrors = 0;

  // NOTE: this supports registering one callback per event, just for keeping it simple
  function on(event: string, callback: unknown) {
    events[event] = callback;
  }

  function commitData() {
    Logger.debug(`Committing Batch [size: ${buffer.length}]`);

    const data = [...buffer];
    buffer = [];

    const insert = bulkInsert(data)
      .then(() => {
        Logger.info('Batch was successfully committed');
        rowsInserted += data.length;
      })
      .catch((error) => {
        bulkInsertErrors++;
        saveToErrorLog(data);
        Logger.error(error);
      }).finally(() => {
        resume();
      });

    promises.push(insert);
  }

  getData((datum) => {
    const result = transformDatum(datum);
    if (result) buffer.push(result);

    if (buffer.length === BATCH_RECORDS) {
      pause();
      commitData();
    }
  });

  onEnd(async () => {
    if (buffer.length) {
      commitData();
    }

    await Promise.all(promises);

    events['finish']({
      duration: `${Date.now() - startTime} ms`,
      rowsInserted,
      errors: {
        transform: transformationErrors,
        bulkInsert: bulkInsertErrors
      }
    });
  });

  return {
    on
  };
}

function transformDatum(datum: string): Matchup | null {
  try {
    if (!datum || datum === '[' || datum === ']') return null;
    const {
      p_match_id,
      champion,
      totaldamagedealttochampions,
      gameversion,
      goldearned,
      win,
      kills,
      assists,
      deaths,
      totalminionskilled,
      summonername,
      visionscore
    } = JSON.parse(datum.replace(/^,*/, ''));
    const [region, , lane] = p_match_id.split('_');

    return {
      champion,
      damage_dealt_to_champions: totaldamagedealttochampions,
      game_version: gameversion,
      gold_earned: goldearned,
      win: win === 'true',
      minions_killed: totalminionskilled,
      kda: (kills + assists) / deaths === 0 ? 1 : deaths,
      lane: lane === 'utility' ? 'support' : lane,
      region,
      summoner_name: summonername,
      vision_score: visionscore
    };
  } catch (error: any) {
    Logger.error(error.message);
    Logger.debug(datum);
    return null;
  }
}