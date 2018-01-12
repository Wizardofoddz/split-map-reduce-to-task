package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/computes/split-map-reduce-to-task/splitmapreducejob"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "split-map-reduce-to-task"
	app.Version = Version
	app.CustomAppHelpTemplate = HelpTemplate
	app.Usage = "USAGE: cat file | ./split-map-reduce-to-task"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ipfs-url, i",
			EnvVar: "IPFS_URL",
			Usage:  "URL to an IPFS API instance",
			Value:  "http://localhost:5001",
		},
	}
	app.Action = run
	app.Run(os.Args)
}

func fatalIfErr(err error, message ...interface{}) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, message...)
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func run(c *cli.Context) {
	fmt.Println("running")
	fi, err := os.Stdin.Stat()
	fatalIfErr(err, "Failed to stat stdin")

	if (fi.Mode() & os.ModeCharDevice) != 0 {
		cli.ShowAppHelp(c)
		fmt.Fprintln(os.Stderr, "\nError, no data found in stdin")
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	fatalIfErr(err, "Failed to read from stdin")

	ipfsURL, err := url.Parse(c.GlobalString("ipfs-url"))
	fatalIfErr(err, "Failed to parse ipfs-url")

	job := splitmapreducejob.New(*ipfsURL)
	err = json.Unmarshal(data, job)
	fatalIfErr(err, "Failed to parse JSON from stdin")

	destinationAddr, err := job.StoreDestination()
	fatalIfErr(err, "Failed to store Destination")

	definitionAddrs, err := job.StoreTaskDefinitions()
	fatalIfErr(err, "Failed to Store Task Definitions")

	fmt.Println("Destination")
	fmt.Println("===========")
	fmt.Println(destinationAddr)
	fmt.Println("")
	fmt.Println("Task Definitions")
	fmt.Println("================")
	for _, addr := range definitionAddrs {
		fmt.Println(addr)
	}

	// taskAddrs, err := job.StoreTasks()
	// fatalIfErr(err, "Failed to generate Task Definitions")

	// output, err := json.MarshalIndent(taskDefinitions, "", "  ")
	// fatalIfErr(err, "Failed to generate output")
	// fmt.Println(string(output))
}
