package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	pgHost     = "localhost"
	pgPort     = 6000
	pgUser     = "postgres"
	pgPassword = "password"
	pgDBname   = "demo"
)

var redisHostList []string = []string{"localhost:7000", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005"}

var rdb *redis.ClusterClient
var db *sql.DB
var ctx context.Context

func handleInsert(w http.ResponseWriter, r *http.Request) {
	log.Info("Inserting into postgres and redis")

	k := fmt.Sprintf("%d", rand.Intn(10000000))
	v := rand.Intn(10000000)

	c := rdb.Set(ctx, k, v, 0)
	if c.Err() != nil {
		fmt.Fprintf(w, "failed to insert to redis: %s", c.Err().Error())
		log.Fatal(c.Err())
	}
	log.Infof("Successfully inserted %s -> %d into redis", k, v)

	_, err := db.Exec("INSERT INTO t(id, value) VALUES ($1, $2)", k, v)
	if err != nil {
		fmt.Fprintf(w, "failed to insert to postgres: %s", err.Error())
		log.Fatal(err)
	}
	log.Infof("Successfully inserted %s -> %d into postgres", k, v)

	fmt.Fprintf(w, "successfully inserted %s -> %d into redis and postgres", k, v)
}

func main() {
	var err error

	rand.Seed(time.Now().UnixNano())

	db, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPassword, pgDBname),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Postgres ping successful")
	defer db.Close()

	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:      redisHostList,
		MaxRetries: 10,
	})

	ctx = context.Background()

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Redis Cluster Ping successful")

	http.HandleFunc("/insert", handleInsert)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
