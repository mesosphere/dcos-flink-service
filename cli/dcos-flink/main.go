// Created by @OhRobin and @joerg84

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"mime/multipart"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mesosphere/dcos-commons/cli"
	"github.com/mesosphere/dcos-commons/cli/client"
	"github.com/mesosphere/dcos-commons/cli/config"
	"gopkg.in/alecthomas/kingpin.v2"
)


func main() {
	app := cli.New()

	handleListJobsSection(app)
	handleJobSection(app)
	handleRunSection(app)
	handleCancelSection(app)
	handleJarsSection(app)
	handleUploadSection(app)

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

//upload
 type UploadHandler struct {
 	filename string
 }
 
 func (cmd *UploadHandler) runUpload(c *kingpin.ParseContext) error {

 	//TODO: x509 auth instead of https to http change
 	url := client.OptionalCLIConfigValue("core.dcos_url") //TODO this should be a RequiredCLIConfigValue
 	url = strings.Replace(url,"https://", "http://", 1)
 	serviceName := config.ServiceName
 	url = fmt.Sprintf("%s/service/%s/jars/upload", url, serviceName)
 	
 	fmt.Println(url)


 	//create multipart payload
  payload := &bytes.Buffer{}
  bodyWriter := multipart.NewWriter(payload)

  fileWriter, err := bodyWriter.CreateFormFile("jarfile", filepath.Base(cmd.filename))
  if err != nil {
      fmt.Println("error writing to buffer")
      return err
  }

   // open file handle
  fh, err := os.Open(cmd.filename)
  if err != nil {
      fmt.Println("error opening file")
      return err
  }

  //iocopy
  _, err = io.Copy(fileWriter, fh)
  if err != nil {
      return err
  }

  // create request
  contentType := bodyWriter.FormDataContentType()
  bodyWriter.Close()

 	req, err := http.NewRequest("POST", url, payload)
  if err != nil {
      return err
  }
  req.Header.Set("Content-Type", contentType)
  req.Header.Add("authorization", fmt.Sprintf("token=%s", client.OptionalCLIConfigValue("core.dcos_acs_token")))

  res, err := http.DefaultClient.Do(req)

  defer res.Body.Close()
  
  // handle response
  resp_body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      return err
  }
  if res.StatusCode != 200 {
  	fmt.Println(res.Status)
  	return errors.New("Upload did not succeed.")
  }
  fmt.Println(string(resp_body))
  return nil
 }
 
 func handleUploadSection(app *kingpin.Application) {
 	cmd := &UploadHandler{}
 	upload := app.Command("upload", "Upload flink jar to run").Action(cmd.runUpload)
 	upload.Arg("jar file", "jar file to upload").Required().StringVar(&cmd.filename)
 }
