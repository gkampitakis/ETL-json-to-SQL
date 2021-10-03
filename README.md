# ETL-json-to-SQL

Load json file, transform it and save it to PostgreSQL


## Description


## Contents

<!-- Dataset  -->
https://www.kaggle.com/jasperan/league-of-legends-1v1-matchups-results?select=matchups.json
https://github.com/vitaly-t/pg-promise/wiki/Data-Imports
https://github.com/vitaly-t/pg-promise/wiki/Performance-Boost
https://stackoverflow.com/questions/37300997/multi-row-insert-with-pg-promise
https://www.postgresql.org/docs/current/populate.html

```sql
\copy matchups(champion,damage_dealt_to_champions,game_version,gold_earned,win,minions_killed,kda,lane,region,summoner_name,vision_score) from '/usr/log.csv' (FORMAT csv,DELIMITER ',');
```

Number of records: 1_312_252

 ETL_ENV=node-solution docker-compose up -d

FILE_PATH=../data-to-load/matchups.json BATCH_RECORDS=8000 yarn start
// NOTE: add in starting message options that the system starts