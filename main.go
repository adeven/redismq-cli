package main

import (
	"fmt"
	"github.com/adeven/redis"
	"os"
	"strconv"
)

var (
	RedisURL      = ""
	RedisHost     = os.Getenv("REDISMQ_HOST")
	RedisPort     = os.Getenv("REDISMQ_PORT")
	RedisDB       = os.Getenv("REDISMQ_DB")
	RedisPassword = os.Getenv("REDISMQ_PASSWORD")

	RedisDBInt = int64(9)
)

func checkRedisConnection() {
	if RedisHost == "" {
		RedisHost = "localhost"
		fmt.Fprintln(os.Stderr, "WARNING: REDISMQ_HOST not found in env, falling back to 'localhost'")
	}
	if RedisPort == "" {
		RedisPort = "6379"
		fmt.Fprintln(os.Stderr, "WARNING: REDISMQ_PORT not found in env, falling back to '6379'")
	}
	if RedisDB == "" {
		RedisDB = "9"
		fmt.Fprintln(os.Stderr, "WARNING: REDISMQ_DB not found in env, falling back to '9'\n")
	}

	RedisURL = RedisHost + ":" + RedisPort

	RedisDBInt, err := strconv.ParseInt(RedisDB, 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	redisClient := redis.NewTCPClient(RedisURL, RedisPassword, RedisDBInt)

	ping := redisClient.Ping()
	if ping.Val() != "PONG" {
		fmt.Fprintf(os.Stderr, "Could not establish connection to %s db %s\n", RedisURL, RedisDB)
		os.Exit(2)
	}
	redisClient.Close()
}

func main() {
	checkRedisConnection()

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	usage()
}
