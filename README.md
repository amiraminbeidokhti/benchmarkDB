# Benchmark Databases

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
#### I)MySQL

Run databases:
```bash
$ docker run --network=my-net --name mySqlServer -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server1/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server1/backup:/backup -h mysql1 mysql:5.7.31 
$ docker run --network=my-net --name mySqlSlave -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server2/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server2/backup:/backup -h mysql2 mysql:5.7.31
```
Config replication MySQL:
```bash
$ docker exec -ti mySqlServer sh -c "mysql -uroot -proot"
$ source /backup/initdb.sql
```
We need File name and Position at the next step.
```bash
$ exit
$ docker exec -ti mySqlSlave sh -c "mysql -uroot -proot"
$ stop slave;
$ CHANGE MASTER TO MASTER_HOST='mySqlServer', MASTER_USER='replicator', MASTER_PASSWORD='replicator', MASTER_LOG_FILE='mysql-bin.000003', MASTER_LOG_POS=154;
$ start slave;
$ exit
```

#### II) PostgreSQL

Run databases:
```bash
$ docker run -d --network=my-net --name postgreSqlServer -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_DATABASE=test -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
$ docker run -d --name postgreSqlSlave --network=my-net -e POSTGRESQL_REPLICATION_MODE=slave -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_MASTER_HOST=postgreSqlServer -e POSTGRESQL_MASTER_PORT_NUMBER=5432 -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
```

#### III) REDIS

Run databases:
```bash
$ docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis
$ docker run --name redisServerReplica -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis
```

#### IV) MYSTORAGE

Run database:
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

We create two instance of different databases. The first instance is master and the second one is slave. Slave replicate all the data of master instanse.