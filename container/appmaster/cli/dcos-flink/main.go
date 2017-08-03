// Created by @OhRobin

package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
)


func main() {
	app := cli.New()

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
	cancel := app.Command("cancel", "Cancel flink job").Action(cmd.runCancel)
	cancel.Arg("job id", "job id of flink").Required().StringVar(&cmd.cancel)
}


type UploadHandler struct {
	upload string
}

func (cmd *UploadHandler) runUpload(c *kingpin.ParseContext) error {
	//TODO x509 auth instead of https to http change
	url := client.OptionalCLIConfigValue("core.dcos_url") //TODO this should be a RequiredCLIConfigValue
	url = strings.Replace(url,"https://", "http://", 1)
	url = fmt.Sprintf("%s/service/flink/jars/upload", url)

	//create multipart payload
	payload := strings.NewReader(fmt.Sprintf("------FormBoundary@OhRobin\r\nContent-Disposition: form-data; name=\"jarfile\"; filename=\"%s\"\r\nContent-Type: application/java-archive\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--", cmd.upload))

	req, _ := http.NewRequest("POST", url, payload)

	//fetch the Auth token from the main CLI.
	req.Header.Add("authorization", fmt.Sprintf("token=%s", client.OptionalCLIConfigValue("core.dcos_acs_token")))
	req.Header.Add("content-type", "multipart/form-data; boundary=----FormBoundary@OhRobin")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
    fmt.Println("Error: %s\n", err)
    return nil
}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return nil
}

func handleUploadSection(app *kingpin.Application) {
	cmd := &UploadHandler{}
	upload := app.Command("upload", "Upload flink jar to run").Action(cmd.runUpload)
	upload.Arg("jar file", "jar file to upload").Required().StringVar(&cmd.upload)
}
