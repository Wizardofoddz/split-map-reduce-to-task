package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	POLYMORPH "github.com/computes/go-ipld-polymorph"
	SPLITMAPREDUCE "github.com/computes/go-sdk/pkg/patterns/splitmapreduce"
	"github.com/urfave/cli"
)

// Version will be replaced by CI
var Version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "split-map-reduce-to-task"
	app.Version = Version
	app.CustomAppHelpTemplate = HelpTemplate
	app.Usage = "USAGE: cat file | ./split-map-reduce-to-task"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ipfs-url, i",
			EnvVar: "COMPUTES_IPFS_URL",
			Usage:  "URL to an IPFS API instance",
			Value:  "http://localhost:5001",
		},
		cli.StringFlag{
			Name:   "http-api-url",
			EnvVar: "COMPUTES_HTTP_API_URL",
			Usage:  "URL to the computes daemon",
			Value:  "http://localhost:8189",
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

	httpAPIURL, err := url.Parse(c.GlobalString("http-api-url"))
	fatalIfErr(err, "Failed to parse http-api-url")

	var jobInput JobInput
	err = json.Unmarshal(data, &jobInput)
	fatalIfErr(err, "Failed to parse JSON from stdin")

	inputPoly := POLYMORPH.New(*ipfsURL)
	err = inputPoly.UnmarshalJSON(jobInput.Input)
	fatalIfErr(err, "Failed to convert input to polymorph")

	job := SPLITMAPREDUCE.New(&SPLITMAPREDUCE.Options{
		IPFSURL:      ipfsURL,
		HTTPAPIURL:   httpAPIURL,
		UUID:         jobInput.UUID,
		SplitInput:   inputPoly,
		SplitRunner:  jobInput.Split.Runner,
		MapRunner:    jobInput.Map.Runner,
		ReduceRunner: jobInput.Reduce.Runner,
	})

	err = job.Create()
	fatalIfErr(err, "Failed to create Job")

	err = job.Run()
	fatalIfErr(err, "Failed to run Job")

	fmt.Println("dataset")
	fmt.Println("=======")
	fmt.Println(job.ResultCID)
	fmt.Println("")
	fmt.Println("task definitions")
	fmt.Println("================")
	fmt.Println(job.SplitTaskDefinitionCID)
	fmt.Println(job.MapTaskDefinitionCID)
	fmt.Println(job.ReduceTaskDefinitionCID)
	fmt.Println("")
	fmt.Println("task")
	fmt.Println("====")
	fmt.Println(job.SplitTaskCID)
}
