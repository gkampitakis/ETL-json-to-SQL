# ETL-json-to-SQL

Load json file, transform it and save it to PostgreSQL

## Description

This repository contains two solutions for "extracting", "transforming" and "loading" data to Postgresql. One solution is written in NodeJS and the second one in Golang.

The focus is how to load data in a performant way (streaming ??) and after transforming save them in Postgresql (in bulk inserts).
### Dataset

The file used for running the ETL pipeline was taken from [here](https://www.kaggle.com/jasperan/league-of-legends-1v1-matchups-results?select=matchups.json). 
If it's not available you can also find it in this [Google Drive](https://drive.google.com/file/d/1DTq50VffBrT4NCKAj2gFVhbGhplHG_6w/view?usp=sharing). The default path that the file is loaded from is `./data-to-load` but you can specify another path by setting the env variable `FILE_PATH`.

> Number of records inside the matchups.json 1.312.252

This is the origin

```json
{
  "p_match_id":"TR1_1201957752_top",
  "goldearned":14425,
  "totalminionskilled":194,
  "win":"false",
  "kills":14,
  "assists":5,
  "deaths":7,
  "champion":"Kassadin",
  "visionscore":17,
  "puuid":"phduyQLB8gBjUerFwiVOtyLLHE9jxw7Jq7dwab_CtRddAvzJ7L1uo5kWzLTKSqStAzml_3yGHiNPFA",
  "totaldamagedealttochampions":33426,
  "summonername":"Borke",
  "gameversion":"11.14.384.6677"
}
```

and we load it in Postgresql as 

```sql
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
```

## Useful commands

__Insert Data to Postgres from CSV:__
```sql
\copy matchups(champion,
damage_dealt_to_champions,game_version,gold_earned,win,minions_killed,kda,lane,region,summoner_name,vision_score) from '/usr/log.csv' (FORMAT csv,DELIMITER ',');
```

__Get number of unique records in table:__

```sql
SELECT COUNT(*) FROM (SELECT DISTINCT * FROM matchups) as unique_rows;
```

__Connect to postgresql container:__

```bash
docker exec -it postgres psql -U ETL_user -d ETL_db
```

## Running the project

For running etl in both versions
- you need to have setup the `.env` or providing the correct environmental variables
- running postgres, in the repo a `docker-compose.yaml` is provided for running postgres. You
 can run it with `make docker-start`

<details>
  <summary>Node</summary>
  
  - Build code `make node-build`
  - Run code (after building) `make node-run`
</details>

## Resources

- [Data Imports](https://github.com/vitaly-t/pg-promise/wiki/Data-Imports) wiki doc for the NodeJS Postgres Driver
- [Performance Boost](https://github.com/vitaly-t/pg-promise/wiki/Performance-Boost) wiki doc for the NodeJS Postgres Driver
- [Multi row insert with pg-promise](https://stackoverflow.com/questions/37300997/multi-row-insert-with-pg-promise)
- [Postgresql - Populating a database](https://www.postgresql.org/docs/current/populate.html)
