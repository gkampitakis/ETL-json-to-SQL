export interface Matchup {
  gold_earned: number;
  minions_killed: number;
  kda: number;
  champion: string;
  vision_score: number;
  summoner_name: string;
  win: boolean,
  game_version: string;
  damage_dealt_to_champions: number;
  lane: string;
  region: string;
}

type Fn<T = void> = () => T;

type Callback<T> = (callback: Fn) => T;

export interface BatchProcessingParams {
  getData: (callback: (data: string) => void) => void;
  pause: Fn<Interface>;
  commitData: (data: Matchup[]) => Promise<null>;
  resume: Fn<Interface>;
}
