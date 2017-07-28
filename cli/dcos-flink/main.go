package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"os"
	"strings"
	"os/exec"
	"mime/multipart"
	"bytes"
	"net/http"
	"io"
	"path/filepath"
	"net/textproto"
)

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

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)


	_ = writer.WriteField("Content-Disposition", "form-data; name='jarfile'; filename='YourFileName.jar'")

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}


func (cmd *UploadHandler) runUpload(c *kingpin.ParseContext) error {
	//TODO remove hardcoding
	//TODO this does't work right now:(
	mh := make(textproto.MIMEHeader)
  mh.Set("Content-Type", "text/plain")
  mh.Set("Content-Disposition", "form-data; name=\"readme\"; filename=\"README.md\"")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("------------------------------printing request-------------------------")
	req, err := newfileUploadRequest("jars/upload" , "jarfile", "/Users/robinoh/Downloads/flink-1.3.1/examples/batch/WordCount.jar")

	fmt.Println(req)
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println("------------------------------printing error-------------------------")
	fmt.Println(err)
	fmt.Println("-----------------------------------------------------------------------")
	body_buf := bytes.NewBufferString("")
	response, err := client.HTTPServicePostData("jars/upload", body_buf, "application/x-java-archive")
	client.PrintJSONBytes(response)
	fmt.Println(err)
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
