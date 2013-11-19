package main

import (
	"fmt"
	"github.com/adeven/redismq"
	"os"
	"text/template"
)

var cmdInfo = &Command{
	Run:   runInfo,
	Usage: "info [queue name]",
	Short: "info returns statistical data for a queue or all queues",
	Long: `
Info returns an overview of stats redismq. Including data about the consumers.
When provided with a queue name only information about this queue is shown.
`,
}

func runInfo(cmd *Command, args []string) {
	if len(args) > 1 {
		cmd.printUsage()
		os.Exit(2)
	}

	ob := redismq.NewObserver(RedisURL, RedisPassword, RedisDBInt)
	queues, err := ob.GetAllQueues()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching queue statistics: %s\n", err)
		os.Exit(2)
	}

	if len(queues) == 0 {
		fmt.Println("no queues found in database")
		return
	}

	if len(args) == 1 {
		name := args[0]
		if !stringInSlice(name, queues) {
			fmt.Fprintf(os.Stderr, "queue with the name %s does not exists\n", name)
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

Queue:	{{$queue}}
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
