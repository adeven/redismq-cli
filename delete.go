package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdDelete = &Command{
	Run:   runDelete,
	Usage: "delete [queue name]",
	Short: "delete a queue",
	Long:  `Delete completely removes a queue from the database incl. all consumers.`,
}

func runDelete(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	queue, err := redismq.SelectQueue(RedisURL, RedisPassword, RedisDBInt, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "queue with the name %s does not exists\n", name)
		os.Exit(2)
	}

	err = queue.Delete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error deleting the queue %s: %s\n", name, err)
		os.Exit(2)
	}
	fmt.Printf("deleted queue with the name %s\n", name)
}
