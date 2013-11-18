package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdList = &Command{
	Run:   runList,
	Usage: "list",
	Short: "lists all queues",
	Long:  `Lists all redismq queues found in this redis database.`,
}

func runList(cmd *Command, args []string) {
	if len(args) > 0 {
		cmd.printUsage()
		os.Exit(2)
	}
	queues, err := redismq.ListQueues(RedisURL, RedisPassword, RedisDBInt)
	if err != nil {
		fmt.Printf("error fetching queue list: %s", err)
	}

	if len(queues) == 0 {
		fmt.Println("No redismq queues found in redis")
		return
	}

	fmt.Println("Found following queues:")
	for _, queue := range queues {
		fmt.Println("\t" + queue)
	}
	fmt.Println("")
	fmt.Println("To show details for a specific queue use:\n\tredismq_cli info [queue name]")
}
