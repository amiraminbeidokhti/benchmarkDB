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
#### I)MySQL:

Run databases:
```bash
$ docker run --network=my-net --name mySqlServer -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server1/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server1/backup:/backup -h mysql1 mysql:5.7.31 
$ docker run --network=my-net --name mySqlServerReplica -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/server2/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/server2/backup:/backup -h mysql2 mysql:5.7.31
```
Config replication MySQL:
```bash
$ docker exec -ti mySqlServer sh -c "mysql -uroot -proot"
$ source /backup/initdb.sql
```
We need File name and Position at the next step.
```bash
$ exit
$ docker exec -ti mySqlServerReplica sh -c "mysql -uroot -proot"
$ stop slave;
$ CHANGE MASTER TO MASTER_HOST='mySqlServer', MASTER_USER='replicator', MASTER_PASSWORD='replicator', MASTER_LOG_FILE='mysql-bin.000003', MASTER_LOG_POS=154;
$ start slave;
$ exit
```

#### II) PostgreSQL:

Run databases:
```bash
$ docker run --name postgreSqlServer --network=my-net -e POSTGRES_PASSWORD=root -d postgres
$ docker run --name postgreSqlServerReplica --network=my-net -e POSTGRES_PASSWORD=root -d postgres
```
Create database and table:
```bash
$ docker exec -it postgreSqlServer psql -U postgres
$ Create database test;
$ \c test
$ CREATE TABLE test (id serial primary key, f1 int, f2 int, f3 int);
$ exit
```
```bash
$ docker exec -it postgreSqlServerReplica psql -U postgres
$ Create database test;
$ \c test
$ CREATE TABLE test (id serial primary key, f1 int, f2 int, f3 int);
$ exit
```

#### III) REDIS:

Run databases:
```bash
$ docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis:latest
$ docker run --name redisServerReplica -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis:latest
```

#### IV) MYSTORAGE:

Run database:
```bash
$ docker run --name redisMyStorageServer --network=my-net -d redis
```

### 4. Run benchmark:
```bash
$ docker build -t benchmark .
$ docker run --network=my-net benchmark
```