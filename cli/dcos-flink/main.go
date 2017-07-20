package main

import (
	"github.com/mesosphere/dcos-commons/cli"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"net/url"
)

func main() {
	app := cli.New()

	cli.HandleDefaultSections(app)
	fmt.Printf("hey")
	kingpin.MustParse(app.Parse(cli.GetArguments()))
}

func (cmd *BrokerHandler) runList(c *kingpin.ParseContext) error {
	cli.PrintJSON(cli.HTTPGet("v1/brokers"))
	return nil
}

func handleTopicSection(app *kingpin.Application) {
	cmd := &TopicHandler{}
	topic.Command(
		"list",
		"Lists all scheduled and running jobs").Action(cmd.runList)
}
