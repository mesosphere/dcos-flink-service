package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"bytes"
	"log"
	"os"
	"os/exec"
)

// "strings"
// "net/http"
// "io/ioutil"
func main() {
	app := cli.New()

	// cli.HandleDefaultSections(app)
	handleListSection(app)
	handleJobSection(app)
	handleRunSection(app)
	handleUploadSection(app)
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
		log.Println(err)
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
	response, err := client.HTTPServicePostQuery(fmt.Sprintf("jars/%s/run", cmd.run), "entry-class=org.apache.flink.examples.java.wordcount.WordCount")
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
		log.Println(err)
	}
	return nil
}

func handleCancelSection(app *kingpin.Application) {
	cmd := &CancelHandler{}
	cancel := app.Command("cancel", "cancel flink job").Action(cmd.runCancel)
	cancel.Arg("job id", "job id of flink").Required().StringVar(&cmd.cancel)
}


type UploadHandler struct {
	upload string
}

func (cmd *UploadHandler) runUpload(c *kingpin.ParseContext) error {
	var out bytes.Buffer
	nodeListExec := exec.Command("curl",
															"--request", "POST",
															"--url", "http://35.163.9.200/service/flink/jars/upload",
															"--header", 	"'authorization:token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDIwMzc0NTgsInVpZCI6ImJvb3RzdHJhcHVzZXIifQ.GFaNZ7li5DIl64tWAWeWmhOzF6VeBubwcV9yKgc9evSxjlEC35JPpn5BpeHVudW54ha6Nd06Gl8YbtkhNK6kHyiSz5OzatfAW_rEApD1orgBRNkZ2N26q7ELpEZwGn1V8NX7OljKM61Lcn5rZFSC6YQXXQTnOtcmq8ntXycNi7xpBLCa3G0n0PNIuAbWevHVbUW8hFU4SELhnaCnRO_SM7F0S1grWpVXX0Xt99CZppb_cgCzbFJi6DBtI5jWsVOVJcDPlvUw1QjnYbwLKwV-jpZKSmCNTqqiJy04bypJKGE9F1ZtZAs0AC8l9YYktrwXDhY93qmgLIS3jAOnB8l5Tw'",
														 	"--header", "'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW'","--header","Expect:",
													 		"--form", "jarfile=@/Users/robinoh/Desktop/dcos-flink-service/container/appmaster/flink/flink-examples/flink-examples-batch/target/WordCount.jar")
	nodeListExec.Stdin = os.Stdin
	nodeListExec.Stdout = &out
	nodeListExec.Stderr = os.Stderr

	err := nodeListExec.Run()
	if err != nil {
		fmt.Printf("[Error] %s\n\n", err)
		fmt.Printf("Unable to run DC/OS command")
		fmt.Printf("Make sure your PATH includes the 'dcos' executable.\n")
	}
}

func handleUploadSection(app *kingpin.Application) {
	cmd := &UploadHandler{}
	upload := app.Command("upload", "Upload flink jar to run").Action(cmd.runUpload)
	upload.Arg("jar file", "jar file to upload").Required().StringVar(&cmd.upload)
}
