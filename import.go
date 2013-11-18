package main

import (
	"bufio"
	"fmt"
	"github.com/adeven/redismq"
	"os"
)

var cmdImport = &Command{
	Run:   runImport,
	Usage: "import [-f filename] [-c count] [-o offset] [queue name]",
	Short: "import each new line as package into a queue",
	Long: `
Imports files to queues. Each line will be a new package.

Options:

    -f		file name to read from
    -c		the number of lines to read
    -o 		the number of lines to skip before reading
`,
}

var (
	fileName string
	offset   int
	maxCount int
)

func init() {
	cmdImport.Flag.StringVar(&fileName, "f", "", "file name")
	cmdImport.Flag.IntVar(&offset, "o", 0, "offset")
	cmdImport.Flag.IntVar(&maxCount, "c", 0, "count")
}

func runImport(cmd *Command, args []string) {
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

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error opening file %s: %s", fileName, err.Error())
		os.Exit(2)
	}
	reader := bufio.NewReader(file)

	lineCount := 0
	imported := 0

	for err == nil {
		lineCount++
		if maxCount != 0 && imported >= maxCount {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if offset != 0 && lineCount > offset {
			queue.Put(line)
			imported++

			fmt.Printf("imported line %d %s", lineCount, line)
		}
	}
	fmt.Printf("finished importing %d package(s)\n", imported)

}