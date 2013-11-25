package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdDelete = &Command{
	Run:   runDelete,
	Usage: "delete [options] [queue name]",
	Short: "delete a queue",
	Long:  `Delete completely removes a queue from the database incl. all consumers.`,
}

func init() {
	cmdDelete.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdDelete.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdDelete.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdDelete.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
}

func runDelete(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	queue, err := redismq.SelectQueue(RedisHost, RedisPort, RedisPassword, RedisDB, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "queue with the name %s does not exists on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
		os.Exit(2)
	}

	err = queue.Delete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error deleting the queue %s on %s:%s db %d: %s\n", name, RedisHost, RedisPort, RedisDB, err)
		os.Exit(2)
	}
	fmt.Printf("deleted queue with the name %s on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
}
