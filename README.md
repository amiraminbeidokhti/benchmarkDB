# Benchmark Databases

In this project we benchmark four different ways of keeping our data for [Gimulator](https://github.com/Gimulator/Gimulator) project.

## Table of Contents

- [Getting Started](#getting-started)
- [Docker compose](#docker-compose)
- [Docker Cli](#docker-cli)
  - [Networks](#create-proper-networks)
  - [Database](#setup-databases)
    - [MySQL](#1.-mysql)
        - [Synchronous](#synchronous)
        - [Asynchronous](#asynchronous)
    - [PostgreSQL](2.-postgresql)
    - [Redis](#3.-reids)
    - [MYSTORAGE](#4.-mystorage)
  - [Benchmark](#benchmark)
- [Description](#description)
- [Contributing](#contributing)

## Getting Started

### Get the project using this command:

```
$ git clone https://github.com/amiraminbeidokhti/benchmarkDB.git
```

## Docker Compose

You can run the project in an **asynchronous** way with this easy command:
```bash
$ docker-compose up
```

## Docker Cli

### Create proper networks:

```bash
$ docker network create -d bridge my-net
$ docker network create cluster --subnet=192.168.0.0/16
```

### Setup databases:
#### 1. MySQL 

#### Synchronous

If you want add other nodes, you must edit mysql-docker/7.6/cnf/mysql-cluster.cnf and specify how many nodes you want. Then you must edit mysql-docker/7.6/cnf/my.cnf and modify the ndb-connectstring to match the ndb_mgmd node.
Build the image:
```bash
$ docker build -t mysql-cluster mySQL/mysql-docker/7.5
```
Creating the manager node:
```bash
$ docker run -d --net=cluster --name=management1 --ip=192.168.0.2 mysql-cluster ndb_mgmd
```
Creating the data nodes:
```bash
$ docker run -d --net=cluster --name=ndb1 --ip=192.168.0.3 mysql-cluster ndbd
$ docker run -d --net=cluster --name=ndb2 --ip=192.168.0.4 mysql-cluster ndbd
```
Creating the mySQL nodes:
```bash
$ docker run -d --net=cluster --name=mysql1 -e MYSQL_DATABASE=test --ip=192.168.0.10 -e MYSQL_ROOT_PASSWORD=root mysql-cluster mysqld
$ docker run -d --net=cluster --name=mysql2 -e MYSQL_DATABASE=test --ip=192.168.0.9 -e MYSQL_ROOT_PASSWORD=root mysql-cluster mysqld
```
Giving the proper privilage to root user:
```bash
$ docker exec -it mysql1 mysql -uroot -proot
$ GRANT ALL ON *.* to root@'%' IDENTIFIED BY 'root';
$ FLUSH PRIVILEGES;
$ exit
```
#### Asynchronous

_FIRST APPROACH_

Run mySQL server and slaves, and config them by environment variables:
```bash
$ docker run --name mySqlServer --network=my-net -d -e MYSQL_ROOT_PASSWORD=root -e MYSQL_REPLICATION_MODE=master -e MYSQL_REPLICATION_USER=my_repl_user -e MYSQL_REPLICATION_PASSWORD=my_repl_password -e MYSQL_USER=my_user -e MYSQL_DATABASE=test -e ALLOW_EMPTY_PASSWORD=yes bitnami/mysql
$ docker run --name mySqlSlave -d --network=my-net -e MYSQL_REPLICATION_MODE=slave -e MYSQL_REPLICATION_USER=my_repl_user -e MYSQL_REPLICATION_PASSWORD=my_repl_password -e MYSQL_MASTER_HOST=mySqlServer bitnami/mysql
```

_SECOND APPROACH_

Run mySQL server:
```bash
$ docker run --network=my-net --name mySqlServer -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/asynchronous/server1/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/asynchronous/server1/backup:/backup -h mysql1 mysql:5.7.31 
```
Run mySQL slave:
```bash
$ docker run --network=my-net --name mySqlSlave -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=test -dit -v $(pwd)/mySQL/asynchronous/server2/conf.d:/etc/mysql/conf.d/ -v $(pwd)/mySQL/asynchronous/server2/backup:/backup -h mysql2 mysql:5.7.31
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

#### 2. PostgreSQL

Run postgreSQL server:

> **_NOTE:_** In order to make replication asynchronously, delete POSTGRESQL_SYNCHRONOUS_COMMIT_MODE=on and POSTGRESQL_NUM_SYNCHRONOUS_REPLICAS=1 from the below command.

```bash
$ docker run -d --network=my-net --name postgreSqlServer -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_DATABASE=test -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password -e POSTGRESQL_SYNCHRONOUS_COMMIT_MODE=on -e POSTGRESQL_NUM_SYNCHRONOUS_REPLICAS=1 bitnami/postgresql
```
Run postgreSQL slave:
```bash
$ docker run -d --name postgreSqlSlave --network=my-net -e POSTGRESQL_REPLICATION_MODE=slave -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_MASTER_HOST=postgreSqlServer -e POSTGRESQL_MASTER_PORT_NUMBER=5432 -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
```

#### 3. REDIS

Run redis server:
```bash
$ docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis
```
Run redis slave:
```bash
$ docker run --name redisSlave -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis
```

#### 4. MYSTORAGE

Run redis slave:
```bash
$ docker run --name redisMyStorageServer --network=my-net -d redis
```

### Benchmark

> **_NOTE:_** In order to make mySQL replication asynchronously, change MYSQL_HOST=192.168.0.10 value to mySqlServer.

> **_NOTE:_** In order to make Redis replication asynchronously, change REDIS_SYNC=1 value to 0 and delete NUM_OF_REPLICA.

Run benchmark:
```bash
$ docker build -t benchmark .
$ docker create --name benchmarkRun -e NUM_OF_DATA=1000 -e LENGTH_OF_DATA=10000 -e REDIS_SYNC=1 -e NUM_OF_REPLICA=1 -e MYSQL_HOST=192.168.0.10 --rm --network=my-net benchmark
$ docker network connect cluster benchmarkRun
$ docker start -a benchmarkRun
```

## Description

We create two instance of each databases. The first instance is master and the second one is slave. Slave replicate all the data of the master instanse. It is useful for high-availability feature.