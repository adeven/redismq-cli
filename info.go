package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
	"text/template"
)

var cmdInfo = &Command{
	Run:   runInfo,
	Usage: "info [options] [queue name]",
	Short: "info returns statistical data for a queue or all queues",
	Long: `
Info returns an overview of stats redismq. Including data about the consumers.
When provided with a queue name only information about this queue is shown.
`,
}

func init() {
	cmdInfo.Flag.StringVar(&RedisHost, "host", "localhost", "redis hostname")
	cmdInfo.Flag.StringVar(&RedisPort, "port", "6379", "redis port")
	cmdInfo.Flag.StringVar(&RedisPassword, "pass", "", "redis password")
	cmdInfo.Flag.Int64Var(&RedisDB, "db", 9, "redis database")
}

func runInfo(cmd *Command, args []string) {
	if len(args) > 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	ob := redismq.NewObserver(RedisHost, RedisPort, RedisPassword, RedisDB)
	queues, err := ob.GetAllQueues()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching queue statistics from %s:%s db %d: %s\n", RedisHost, RedisPort, RedisDB, err)
		os.Exit(2)
	}

	if len(queues) == 0 {
		fmt.Printf("no queues found in database on %s:%s db %d\n", RedisHost, RedisPort, RedisDB)
		return
	}

	if len(args) == 1 {
		name := args[0]
		if !stringInSlice(name, queues) {
			fmt.Fprintf(os.Stderr, "queue with the name %s does not exists on %s:%s db %d\n", name, RedisHost, RedisPort, RedisDB)
			os.Exit(2)
		}
		ob.UpdateQueueStats(name)
	} else {
		ob.UpdateAllStats()
	}
	statsTemplate.Execute(os.Stdout, struct {
		Queues map[string]*redismq.QueueStat
	}{
		ob.Stats,
	})
}

var statsTemplate = template.Must(template.New("stats").Parse(
	`{{range $queue, $stats := .Queues}}

Queue(s):	{{$queue}}
_____________________________________________________________________
InputRates:	sec		min		hour
		{{$stats.InputRateSecond}}		{{$stats.InputRateMinute}}		{{$stats.InputRateHour}}
WorkRates:	sec		min		hour
		{{$stats.WorkRateSecond}}		{{$stats.WorkRateMinute}}		{{$stats.WorkRateHour}}

AvgInputSize:	sec		min		hour
		{{$stats.InputSizeSecond}}		{{$stats.InputSizeMinute}}		{{$stats.InputSizeHour}}
AvgFailedSize:	sec		min		hour
		{{$stats.FailSizeSecond}}		{{$stats.FailSizeMinute}}		{{$stats.FailSizeHour}}

Consumers:
_____________________________________________________________________
{{range $consumer, $cstats := $stats.ConsumerStats}}
	Consumer:	{{$consumer}}
	_____________________________________________________________
	WorkRates:	sec		min		hour
			{{$cstats.WorkRateSecond}}		{{$cstats.WorkRateMinute}}		{{$cstats.WorkRateHour}}
	=============================================================
{{end}}{{end}}

`))
