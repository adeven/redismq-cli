package main

import (
	"fmt"
	"github.com/adjust/redismq"
	"os"
)

var cmdList = &Command{
	Run:   runList,
	Usage: "list [options]",
	Short: "lists all queues",
	Long:  `Lists all redismq queues found in this redis database.`,
}

func init() {
	cmdList.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdList.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdList.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdList.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
}

func runList(cmd *Command, args []string) {
	if len(args) > 0 {
		cmd.printUsage()
		os.Exit(2)
	}
	ob := redismq.NewObserver(RedisHost, RedisPort, RedisPassword, RedisDB)
	queues, err := ob.GetAllQueues()
	if err != nil {
		fmt.Printf("error fetching queue list: %s", err)
	}

	if len(queues) == 0 {
		fmt.Printf("No redismq queues found in redis on %s:%s db %d\n", RedisHost, RedisPort, RedisDB)
		return
	}

	fmt.Printf("Found following queues on %s:%s db %d:\n", RedisHost, RedisPort, RedisDB)
	for _, queue := range queues {
		fmt.Println("\t" + queue)
	}
	fmt.Println("")
	fmt.Println("To show details for a specific queue use:\n\tredismq-cli info [queue name]")
}
