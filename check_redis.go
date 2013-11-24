package main

import (
	"fmt"
	"github.com/adeven/redis"
	"os"
)

var cmdCheckRedis = &Command{
	Run: runCheckRedis,
}

func init() {
	cmdCheckRedis.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdCheckRedis.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdCheckRedis.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdCheckRedis.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
}

func runCheckRedis(cmd *Command, args []string) {
	redisURL := RedisHost + ":" + RedisPort
	redisClient := redis.NewTCPClient(redisURL, RedisPassword, RedisDB)
	ping := redisClient.Ping()
	if ping.Val() != "PONG" {
		fmt.Fprintf(os.Stderr, "Could not establish connection to %s db %d\n", redisURL, RedisDB)
		os.Exit(2)
	}
	redisClient.Close()
}
