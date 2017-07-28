package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"os"
	"strings"
	"net/http"
	"os/exec"
	"io/ioutil"
)

//	"net/http/httputil"

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


type UploadHandler struct {
	upload string
}

func (cmd *UploadHandler) runUpload(c *kingpin.ParseContext) error {

	url := "http://54.71.5.146/service/flink/jars/upload"

	payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"jarfile\"; filename=\"WordCount.jar\"\r\nContent-Type: application/java-archive\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"filename\"\r\n\r\nWordCount.jar\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authorization", "token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDE2OTEyMTUsInVpZCI6ImJvb3RzdHJhcHVzZXIifQ.DbIxHsISFsKDmLyWuhzVnx7IlfCyZnJ6x86XH6NC4HiuB3BcU8bvJnfrvq5NX9-eJeKfXQ2bAZjaDYQ9pQBTipyxghF6rxMPQ3HqYEftU07ciwwVAAtZHMg56hI5MNY99vyXjdKPG44nbNxe0a_CirfOKEI-ItnCTDCGvG_OJgsbAUESOOHOT0RGyXzoMZpsiam_u8aFgtbfbyScTzKfFjA8C-aRVIk5D-tXSce_AyDNrsGHNVzjxAxhZ_1EZduCMqwMUgP7Si6sw_-jU_xURAJ9bZBgx1K_SCy005HsW3zHgBhMDYDeeFmnx6kRtVdbIa6x2MZGFTsbITpc4dWKBQ")
	req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "c3f685a8-d4dd-05a3-28e9-bc63ff3dfb80")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}

func handleUploadSection(app *kingpin.Application) {
	cmd := &UploadHandler{}
	// job := app.Command("run", "Run flink job")
	app.Command("upload", "Upload flink jar to run").Action(cmd.runUpload)
}


//TODO remove if not used
func runDcosCommand(arg ...string) {
	cmd := exec.Command("dcos", arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[Error] %s\n\n", err)
		fmt.Printf("Unable to run DC/OS command: %s\n", strings.Join(arg, " "))
		fmt.Printf("Make sure your PATH includes the 'dcos' executable.\n")
	}
}
