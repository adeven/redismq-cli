package main

import (
	"fmt"
	"github.com/adjust/redismq"
	"os"
)

var cmdCreate = &Command{
	Run:   runCreate,
	Usage: "create [options] [queue name]",
	Short: "create a new queue",
	Long:  `Create creates a new redismq queue.`,
}

func init() {
	cmdCreate.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdCreate.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdCreate.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdCreate.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
}

func runCreate(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	ob := redismq.NewObserver(RedisHost, RedisPort, RedisPassword, RedisDB)
	queues, _ := ob.GetAllQueues()
	if stringInSlice(name, queues) {
		fmt.Fprintf(os.Stderr, "queue with the name %s already exists on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
		os.Exit(2)
	}
	redismq.CreateQueue(RedisHost, RedisPort, RedisPassword, RedisDB, name)
	fmt.Printf("created queue with the name %s on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
}
