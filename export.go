package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdExport = &Command{
	Run:   runExport,
	Usage: "export [-f filename] [-c count] [-o offset] [-r] [-v] [queue name]",
	Short: "export each package from a queue to a file",
	Long: `
"Exports each package from a queue to a file.
Fetches queue size upon start and then extracts exactly this number of packages.
It's not adviced to read from the queue while exporting when using the requeue flag.

Options:

    -f		file name to write to
    -c		the number of packages to write
    -r		requeue packages (if set the queue will just be rotated)
    -v		verbose mode (display each export)
`,
}

var (
	flagRequeue bool
)

func init() {
	cmdExport.Flag.StringVar(&fileName, "f", "", "file name")
	cmdExport.Flag.Int64Var(&maxCount, "c", 0, "count")
	cmdExport.Flag.BoolVar(&flagVerbose, "v", false, "verbose mode")
	cmdExport.Flag.BoolVar(&flagRequeue, "r", false, "keep packages in queue")
}

func runExport(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	queue, err := redismq.SelectQueue(RedisURL, RedisPassword, RedisDBInt, name)
	if err != nil {
		fmt.Printf("queue with the name %s doesn't exists\n", name)
		os.Exit(2)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("error creating file %s: %s\n", fileName, err.Error())
		os.Exit(2)
	}
	defer file.Close()

	lineCount := int64(0)
	exported := int64(0)
	todo := queue.GetInputLength()

	consumer, err := queue.AddConsumer("redismq_util export command")
	if err != nil {
		fmt.Printf("error registering consumer: %s\n", err.Error())
		os.Exit(2)
	}

	if consumer.HasUnacked() {
		err := consumer.RequeueWorking()
		if err != nil {
			fmt.Printf("error requeing unacked packages %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Println("found unacked packages...requeued")
	}

	for err == nil && exported < todo {
		lineCount++
		if maxCount != 0 && exported >= maxCount {
			break
		}

		p, err := consumer.NoWaitGet()
		if err != nil {
			fmt.Printf("error fetching package from queue: %s", err.Error())
			os.Exit(2)
		}

		if p == nil {
			break
		}

		file.WriteString(p.Payload)

		if flagRequeue {
			err = p.Requeue()
		} else {
			err = p.Ack()
		}
		if err != nil {
			fmt.Printf("error ack/requeing package to queue: %s", err.Error())
			os.Exit(2)
		}

		exported++
		if flagVerbose {
			fmt.Printf("exported package %d %s", lineCount, p.Payload)
		}

	}
	file.Sync()
	fmt.Printf("finished exported %d package(s)\n", exported)
}
