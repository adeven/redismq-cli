package main

import (
	"fmt"
	"os"
)

var (
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int64
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	//cmdCheckRedis.Run(cmdCheckRedis, cmdCheckRedis.Flag.Args())

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
