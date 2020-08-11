FROM golang

ENV MYSQL_USER=root \
    MYSQL_PASSWORD=root \
    MYSQL_HOST=mySqlServer \
    MYSQL_PORT=3306 \
    MYSQL_DBNAME=test \
    # PostgreSql
    POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=root \
    POSTGRES_HOST=postgreSqlServer \
    POSTGRES_PORT=5432 \
    POSTGRES_DBNAME=test \
    # Redis
    REDIS_HOST=redisServer \
    REDIS_PORT=6379 \
    # MyStorage
    REDIS_HOST_MYSTORAGE=redisMyStorageServer \
    REDIS_PORT_MYSTORAGE=6379 \
    # DATA
    NUM_OF_DATA=1000 \
    LENGTH_OF_DATA=10000


WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD [ "go", "test", "-bench=." ]

