package main

import (
	"bufio"
	"fmt"
	"github.com/adjust/redismq"
	"os"
)

var cmdImport = &Command{
	Run:   runImport,
	Usage: "import [options] -f filename [-c count] [-o offset] [-v] [queue name]",
	Short: "import each new line as package into a queue",
	Long: `
Imports files to queues. Each line will be a new package.

Options:

    -f		file name to read from
    -c		the number of lines to read
    -o 		the number of lines to skip before reading
    -v		verbose mode (display each import)
`,
}

var (
	fileName    string
	offset      int64
	maxCount    int64
	flagVerbose bool
)

func init() {
	cmdImport.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdImport.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdImport.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdImport.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
	cmdImport.Flag.StringVar(&fileName, "f", "", "file name")
	cmdImport.Flag.Int64Var(&offset, "o", 0, "offset")
	cmdImport.Flag.Int64Var(&maxCount, "c", 0, "count")
	cmdImport.Flag.BoolVar(&flagVerbose, "v", false, "verbose mode")
}

func runImport(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	name := args[0]
	queue, err := redismq.SelectQueue(RedisHost, RedisPort, RedisPassword, RedisDB, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "queue with the name %s doesn't exists on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
		os.Exit(2)
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %s: %s", fileName, err.Error())
		os.Exit(2)
	}
	defer file.Close()
	// TODO fall back to stdin
	reader := bufio.NewReader(file)

	lineCount := int64(0)
	imported := int64(0)

	for err == nil {
		lineCount++
		if maxCount != 0 && imported >= maxCount {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if line == "\n" {
			if flagVerbose {
				fmt.Printf("skipping empty line %d\n", lineCount)
			}
			continue
		}

		if lineCount > offset {
			queue.Put(line)
			imported++
			if flagVerbose {
				fmt.Printf("imported line %d %s", lineCount, line)
			}
		}
	}
	fmt.Printf("finished importing %d package(s) to %s:%s db %d\n", imported, RedisHost, RedisPort, RedisDB)
}
