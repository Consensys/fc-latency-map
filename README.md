# fc-latency-map




## Requirements
* sqlite3
```bash
sudo apt install sqlite3
```

## SQLite commands
1. Open database
```bash
sqlite3 data/database.db

SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite>
```

2. View tables
```bash
sqlite> .tables

>>> miners
```

3. Select miners
```bash
sqlite> select * from miners;

>>> 1|2021-09-07 17:04:53.41480159+02:00|2021-09-07 17:04:53.41480159+02:00||dummyAddress|dummyIp
```

4. Close database
```bash
sqlite> ^C^C^C
```


### Configuration:

Must have an ".env" file with the configurations to run.

Is available an example of configuration on [.env.example](./manager/.env.example)