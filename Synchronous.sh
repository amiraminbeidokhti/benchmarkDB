# Create proper networks:
docker network create -d bridge my-net
docker network create cluster --subnet=192.168.0.0/16



# --MYSQL--
# If you want add other nodes, you must edit mysql-docker/7.6/cnf/mysql-cluster.cnf and specify how many nodes you want.
# Then you must edit mysql-docker/7.6/cnf/my.cnf and modify the ndb-connectstring to match the ndb_mgmd node.
# Build the image:
docker build -t mysql-cluster mySQL/mysql-docker/7.5

# Creating the manager node:
docker run -d --net=cluster --name=management1 --ip=192.168.0.2 mysql-cluster ndb_mgmd

# Creating the data nodes:
docker run -d --net=cluster --name=ndb1 --ip=192.168.0.3 mysql-cluster ndbd
docker run -d --net=cluster --name=ndb2 --ip=192.168.0.4 mysql-cluster ndbd

# Creating the mySQL nodes:
docker run -d --net=cluster --name=mysql1 -e MYSQL_DATABASE=test --ip=192.168.0.10 -e MYSQL_ROOT_PASSWORD=root mysql-cluster mysqld
docker run -d --net=cluster --name=mysql2 -e MYSQL_DATABASE=test --ip=192.168.0.9 -e MYSQL_ROOT_PASSWORD=root mysql-cluster mysqld
###########################################
# Now you have to config your mysql
###########################################
# --PostgreSQL--
docker run -d --network=my-net --name postgreSqlServer -e POSTGRESQL_REPLICATION_MODE=master -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_DATABASE=test -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password -e POSTGRESQL_SYNCHRONOUS_COMMIT_MODE=on -e POSTGRESQL_NUM_SYNCHRONOUS_REPLICAS=1 bitnami/postgresql

# Slave:
docker run -d --name postgreSqlSlave --network=my-net -e POSTGRESQL_REPLICATION_MODE=slave -e POSTGRESQL_USERNAME=postgres -e POSTGRESQL_PASSWORD=root -e POSTGRESQL_MASTER_HOST=postgreSqlServer -e POSTGRESQL_MASTER_PORT_NUMBER=5432 -e POSTGRESQL_REPLICATION_USER=my_repl_user -e POSTGRESQL_REPLICATION_PASSWORD=my_repl_password bitnami/postgresql

# --Redis--
docker run --name redisServer --network=my-net -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_REPLICATION_MODE=master -d bitnami/redis

# Slave:
docker run --name redisSlave -e ALLOW_EMPTY_PASSWORD=yes --network=my-net -e REDIS_REPLICATION_MODE=slave  -d -e REDIS_MASTER_HOST=redisServer -e REDIS_MASTER_PORT_NUMBER=6379 bitnami/redis

# --MyStorage--
docker run --name redisMyStorageServer --network=my-net -d redis



# Benchmark
docker build -t benchmark .
docker create --name benchmarkRun -e NUM_OF_DATA=1000 -e LENGTH_OF_DATA=10000 -e REDIS_SYNC=1 -e NUM_OF_REPLICA=1 -e MYSQL_HOST=192.168.0.10 --rm --network=my-net benchmark
docker network connect cluster benchmarkRun
docker start -a benchmarkRun