package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
)


func main() {
	app := cli.New()
	// cli.HandleDefaultSections(app)
	handleListSection(app)
	handleJobSection(app)
	handleRunSection(app)
	handleCancelSection(app)

	kingpin.MustParse(app.Parse(cli.GetArguments()))
}

func handleListSection(app *kingpin.Application) {
	app.Command("list", "List completed and running jobs").Action(runList)
}

func runList(c *kingpin.ParseContext) error {
	response, err := client.HTTPServiceGet("jobs")
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		fmt.Println(err)
	}
	return nil
}


type InfoHandler struct {
	info string
}

func (cmd *InfoHandler) runInfo(c *kingpin.ParseContext) error {
	var response []byte
	var err error

	if cmd.info == "" {
		response, err = client.HTTPServiceGet("joboverview")
	} else {
		response, err = client.HTTPServiceGet(fmt.Sprintf("jobs/%s", cmd.info))
	}

	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		fmt.Println(err)
	}

	return nil
}


func handleJobSection(app *kingpin.Application) {
	cmd := &InfoHandler{}
	job := app.Command("info", "Summary of Job status.").Action(cmd.runInfo)
	job.Arg("jobid",
					"Summary of one job").StringVar(&cmd.info)
}

type RunHandler struct {
	run string
}

func (cmd *RunHandler) runRun(c *kingpin.ParseContext) error {
	response, err := client.HTTPServicePost(fmt.Sprintf("jars/%s/run", cmd.run))
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		fmt.Println(err)
	}
	return nil
}

func handleRunSection(app *kingpin.Application) {
	cmd := &RunHandler{}
	run := app.Command("run", "Run flink job").Action(cmd.runRun)
	run.Arg("JarID", "The filename provided after uploading Jar file").Required().StringVar(&cmd.run)
}

type CancelHandler struct {
	cancel string
}

func (cmd *CancelHandler) runCancel(c *kingpin.ParseContext) error {
	response, err := client.HTTPServiceDelete(fmt.Sprintf("jobs/%s/cancel", cmd.cancel))
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		fmt.Println(err)
	}
	return nil
}

func handleCancelSection(app *kingpin.Application) {
	cmd := &CancelHandler{}
	cancel := app.Command("cancel", "cancel flink job").Action(cmd.runCancel)
	cancel.Arg("job id", "job id of flink").Required().StringVar(&cmd.cancel)
}
