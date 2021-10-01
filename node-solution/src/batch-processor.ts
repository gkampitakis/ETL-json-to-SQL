import { Matchup, Fn, BatchProcessingParams } from './types';
import Logger from './logger';

const MAX_RECORDS = parseInt(process.env.MAX_RECORDS ?? '5000');

export function batchProcessing({
  getData,
  pause,
  commitData,
  resume
}: BatchProcessingParams,) {
  let buffer: Matchup[] = [];
  // let mux = false;

  getData((data) => {
    const result = transformDatum(data);
    if (result) buffer.push(result);

    if (buffer.length >= MAX_RECORDS) {
      pause();

      Logger.debug(`Committing Batch [size: ${buffer.length}]`);

      commitData(buffer)
        .then(() => {
          Logger.debug('Batch was successfully committed');
          buffer = [];
          resume();
        })
        .catch((error) => {
          Logger.error(error);
        });

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
}

function transformDatum(datum: string): Matchup | null {
  if (datum === '[' || datum === ']') return null;
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
    kda: (kills + assists) / deaths,
    lane: lane === 'utility' ? 'support' : lane,
    region,
    summoner_name: summonername,
    vision_score: visionscore
  };
}