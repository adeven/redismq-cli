package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdInfo = &Command{
	Run:   runInfo,
	Usage: "info [queue name]",
	Short: "info returns statistical data for a queue",
	Long:  `Info returns an overview of stats for a queue. Including data about the consumers.`,
}

func runInfo(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	_, err := redismq.SelectQueue(RedisURL, RedisPassword, RedisDBInt, name)
	if err != nil {
		fmt.Printf("queue with the name %s does not exists\n", name)
		os.Exit(2)
	}

}
