// Created by @OhRobin

package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"log"
)


func main() {
	app := cli.New()

	handleListJobsSection(app)
	handleJobSection(app)
	handleRunSection(app)
	handleCancelSection(app)
	handleJarsSection(app)

	kingpin.MustParse(app.Parse(cli.GetArguments()))
}

//list jobs
func handleListJobsSection(app *kingpin.Application) {
	app.Command("list", "List completed and running jobs").Action(runListJobs)
}

func runListJobs(c *kingpin.ParseContext) error {
	response, err := client.HTTPServiceGet("jobs")
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		log.Println(err)
	}
	return nil
}

//list jars
func handleJarsSection(app *kingpin.Application) {
	app.Command("jars", "List uploaded jar files and associated jar ids").Action(runJars)
}

func runJars(c *kingpin.ParseContext) error {
	response, err := client.HTTPServiceGet("jars")
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		log.Println(err)
	}
	return nil
}


//job info
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
	job := app.Command("info", "Summary of job status").Action(cmd.runInfo)
	job.Arg("job id", "Summary of one job").StringVar(&cmd.info)
}


//run
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
	run.Arg("jar id", "The filename provided after uploading jar file").Required().StringVar(&cmd.run)
}

//cancel job
type CancelHandler struct {
	cancel string
}

func (cmd *CancelHandler) runCancel(c *kingpin.ParseContext) error {
	response, err := client.HTTPServiceDelete(fmt.Sprintf("jobs/%s/cancel", cmd.cancel))
	if err == nil {
		client.PrintJSONBytes(response)
	} else {
		log.Println(err)
	}
	return nil
}

func handleCancelSection(app *kingpin.Application) {
	cmd := &CancelHandler{}
	cancel := app.Command("cancel", "Cancel flink job").Action(cmd.runCancel)
	cancel.Arg("job id", "job id of flink").Required().StringVar(&cmd.cancel)
}
