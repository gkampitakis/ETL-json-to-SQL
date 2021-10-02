CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS matchups (
  id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  gold_earned INTEGER,
  minions_killed INTEGER,
  kda INTEGER,
  champion VARCHAR,
  vision_score INTEGER,
  summoner_name VARCHAR,
  win BOOLEAN,
  game_version VARCHAR,
  damage_dealt_to_champions INTEGER,
  lane VARCHAR,
  region VARCHAR
);