package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdCreate = &Command{
	Run:   runCreate,
	Usage: "create [queue name]",
	Short: "create a new queue",
	Long:  `Create creates a new redismq queue.`,
}

func runCreate(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	queues, _ := redismq.ListQueues(RedisURL, RedisPassword, RedisDBInt)
	if stringInSlice(name, queues) {
		fmt.Printf("queue with the name %s already exists\n", name)
		os.Exit(2)
	}
	redismq.CreateQueue(RedisURL, RedisPassword, RedisDBInt, name)
	fmt.Printf("created queue with the name %s\n", name)
}
