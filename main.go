package main

import (
	"flag"
	"fmt"
	"github.com/adeven/redis"
	"os"
	"strconv"
	"strings"
)

var (
	RedisURL      = os.Getenv("REDISMQ_URL")
	RedisDB       = os.Getenv("REDISMQ_DB")
	RedisPassword = os.Getenv("REDISMQ_PASSWORD")

	RedisDBInt = int64(9)
)

var commands = []*Command{
	cmdCreate,
}

type Command struct {
	// args does not include the command name
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage string // first word is the command name
	Short string // `redismq help` output
	Long  string // `redismq help cmd` output
}

func (c *Command) printUsage() {
	if c.Runnable() {
		fmt.Printf("Usage: redismq %s\n\n", c.Usage)
	}
	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func checkRedisConnection() {
	if RedisURL == "" {
		RedisURL = "localhost:6379"
		//fmt.Fprintln(os.Stderr, "REDISMQ_URL not found in env, falling back to 'localhost:6379'\n")
	}
	if RedisDB == "" {
		RedisDB = "9"
		//fmt.Fprintln(os.Stderr, "REDISMQ_DB not found in env, falling back to '9'\n")
	}
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

}
