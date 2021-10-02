import {
  Matchup,
  BatchProcessingParams,
  BatchProcessingEventCallback
} from './types';
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
  let buffer: Matchup[] = [];
  // let mux = false;

  // NOTE: this supports registering one callback per event, just for keeping it simple
  function on(event: string, callback: unknown) {
    events[event] = callback;
  }

  function commitData(final = false) {
    Logger.debug(`Committing Batch [size: ${buffer.length}]`);

    const data = [...buffer];
    buffer = [];

    return bulkInsert(data)
      .then(() => {
        Logger.debug('Batch was successfully committed');
      })
      .catch((error) => {
        Logger.error(error);
      })
      .finally(() => {
        resume();
      });
  }

  getData((datum) => {
    const result = transformDatum(datum);
    if (result) buffer.push(result);

    if (buffer.length >= BATCH_RECORDS) {
      commitData();
      pause();
      // if (!mux) {
      //   mux = true;
      //   setTimeout(() => {
      //     console.log(buffer);
      //     buffer = [];
      //     mux = false;
      //     resume();
      //   }, 5000);
      // }
    }
  });

  onEnd(async () => {
    // NOTE: If there are some data still buffered
    if (buffer.length) {
      await commitData();
    }

    events['finish']({
      data: 'hello World'
    });
  });

  return {
    on
  };
}

function transformDatum(datum: string): Matchup | null {
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
}