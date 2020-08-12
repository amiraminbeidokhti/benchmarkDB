# Benchmark Databases

In this project we benchmark four different ways of keeping our data.

## Getting Started

### 1. Get the project using this command:

```
$ git clone https://github.com/amiraminbeidokhti/benchmarkDB.git
```

### 2. Create bridge network:

```bash
$ docker network create -d bridge my-net
```

### 3. Setup databases:
#### I) MySQL

Run mySQL server:
```bash
$ docker run --network=my-net --name mySqlServer -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server1/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server1/backup:/backup -h mysql1 mysql:5.7.31 
```
Run mySQL slave:
```bash
$ docker run --network=my-net --name mySqlSlave -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server2/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server2/backup:/backup -h mysql2 mysql:5.7.31
```

Config mySQL slave:
```bash
$ docker exec -it mySqlServer sh -c "mysql -uroot -proot"
$ source /backup/initdb.sql
```
We need the File name and Position that are shown at the terminal for the next step.
```bash
$ exit
$ docker exec -it mySqlSlave sh -c "mysql -uroot -proot"
$ stop slave;
```
Now for the MASTER_LOG_FILE and MASTER_LOG_POS insert the information we gathered at the previous step
```bash
$ CHANGE MASTER TO MASTER_HOST='mySqlServer', MASTER_USER='replicator', MASTER_PASSWORD='replicator', MASTER_LOG_FILE='mysql-bin.000003', MASTER_LOG_POS=154;
$ start slave;
$ exit
```

#### II) PostgreSQL

Run postgreSQL server:
If you want synchronous replication, add -e POSTGRESQL_SYNCHRONOUS_COMMIT_MODE=on and -e POSTGRESQL_NUM_SYNCHRONOUS_REPLICAS=1 to the below command.
```bash
$ docker run -d --network=my-net --name postgreSqlServer -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_DATABASE=test -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
```
Run postgreSQL slave:
```bash
$ docker run -d --name postgreSqlSlave --network=my-net -e POSTGRESQL_REPLICATION_MODE=slave -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_MASTER_HOST=postgreSqlServer -e POSTGRESQL_MASTER_PORT_NUMBER=5432 -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
```

#### III) REDIS

Run redis server:
```bash
$ docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis
```
Run redis slave:
```bash
$ docker run --name redisSlave -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis
```

#### IV) MYSTORAGE

Run redis slave:
```bash
$ docker run --name redisMyStorageServer --network=my-net -d redis
```

### 4. Benchmark

Run benchmark:
```bash
$ docker build -t benchmark .
$ docker run --network=my-net benchmark
```

## Description

We create two instance of each databases. The first instance is master and the second one is slave. Slave replicate all the data of the master instanse. It is useful for high-availability feature.