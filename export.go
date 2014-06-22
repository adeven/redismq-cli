package main

import (
	"fmt"
	"github.com/adjust/redismq"
	"os"
)

var cmdExport = &Command{
	Run:   runExport,
	Usage: "export [options] [-f filename] [-c count] [-o offset] [-r] [-v] [queue name]",
	Short: "export each package from a queue to a file",
	Long: `
"Exports each package from a queue to a file.
Fetches queue size upon start and then extracts exactly this number of packages.
It's not adviced to read from the queue while exporting when using the requeue flag.

Options:

    -f		file name to export to
    -c		the number of packages to export
    -r		requeue packages after export
    -v		verbose mode (print each payload)

If no file name is given the output will be directed to stdout.
`,
}

var (
	flagRequeue bool
)

func init() {
	cmdExport.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdExport.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdExport.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdExport.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
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
	queue, err := redismq.SelectQueue(RedisHost, RedisPort, RedisPassword, RedisDB, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "queue with the name %s does not exists on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
		os.Exit(2)
	}

	var file *os.File
	if fileName != "" {
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating file %s: %s\n", fileName, err.Error())
			os.Exit(2)
		}
		defer file.Close()
	}

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
			fmt.Fprintf(os.Stderr, "error requeing unacked packages %s\n", err.Error())
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
			fmt.Fprintf(os.Stderr, "error fetching package from queue: %s", err.Error())
			os.Exit(2)
		}

		if p == nil {
			break
		}

		if file != nil {
			file.WriteString(p.Payload)
		} else {
			fmt.Print(p.Payload)
		}

		if flagRequeue {
			err = p.Requeue()
		} else {
			err = p.Ack()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error ack/requeing package to queue: %s", err.Error())
			os.Exit(2)
		}

		exported++
		if flagVerbose {
			fmt.Printf("exported package %d %s", lineCount, p.Payload)
		}

	}
	if file != nil {
		file.Sync()
	}
	fmt.Printf("\n\nfinished exporting %d package(s) from %s:%s db %d\n", exported, RedisHost, RedisPort, RedisDB)
}
