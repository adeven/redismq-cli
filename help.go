package main

import (
	"os"
	"text/template"
)

var usageTemplate = template.Must(template.New("usage").Parse(`
Usage: redismq [command] [options] [arguments]

Commands:
{{range .Commands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf "%-8s"}}  {{.Short}}{{end}}{{end}}{{end}}

Run 'redismq help [command]' for details.

`[1:]))

func usage() {
	printUsage()
	os.Exit(2)
}

func printUsage() {
	usageTemplate.Execute(os.Stdout, struct {
		Commands []*Command
	}{
		commands,
	})
}
