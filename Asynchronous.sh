# Create proper networks:
docker network create -d bridge my-net

#--MySQL
docker run --name mySqlServer --network=my-net -d -e MYSQL_ROOT_PASSWORD=root -e MYSQL_REPLICATION_MODE=master -e MYSQL_REPLICATION_USER=my_repl_user -e MYSQL_REPLICATION_PASSWORD=my_repl_password -e MYSQL_USER=my_user -e MYSQL_DATABASE=test -e ALLOW_EMPTY_PASSWORD=yes bitnami/mysql
docker run --name mySqlSlave -d --network=my-net -e MYSQL_REPLICATION_MODE=slave -e MYSQL_REPLICATION_USER=my_repl_user -e MYSQL_REPLICATION_PASSWORD=my_repl_password -e MYSQL_MASTER_HOST=mySqlServer bitnami/mysql

#--PostgreSQL--
docker run -d --network=my-net --name postgreSqlServer -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_DATABASE=test -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql
docker run -d --name postgreSqlSlave --network=my-net -e POSTGRESQL_REPLICATION_MODE=slave -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_MASTER_HOST=postgreSqlServer -e POSTGRESQL_MASTER_PORT_NUMBER=5432 -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql

#--Redis--
docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis
docker run --name redisSlave -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis

#--MyStorage--
docker run --name redisMyStorageServer --network=my-net -d redis


#Benchmark
docker build -t benchmark .
docker create --name benchmarkRun -e NUM_OF_DATA=1000 -e LENGTH_OF_DATA=10000 -e REDIS_SYNC=0 -e MYSQL_HOST=mySqlServer --rm --network=my-net benchmark
docker start -a benchmarkRun