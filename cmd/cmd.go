package cmd

import (
	"context"
	"os"
	"strings"

	"Holmes/run"
	"Holmes/utils"

	"github.com/jawher/mow.cli"
	"github.com/pterm/pterm"
)

func Execute() {
	// create an app
	app := cli.App("Holmes", "FingerPrint Recognition")

	var (
		inputurl   	= app.StringOpt("u inputurl", "", "iput a url to scan")
		inputfile   = app.StringOpt("uf inputfile", "", "iput a urlfile to scan")
		threads   	= app.IntOpt("th threads", 10, "thread num")
		timeout     = app.IntOpt("to timeout", 10, "Request timeout")
		rulefile	= app.StringOpt("rf rule", "","yamlfile path")
	)

	app.Version("v version", "Holmes 1.0.0")
	app.Spec = "-v | (-u=<inputurl> | --uf=<inputfile>) (--rf=<rulefile>)... [--th=<threads>] [--to=<timeout>]"

	app.Action = func() {
		urlSlice := make([]string, 0)
		if *inputurl != "" {
			urlSlice = append(urlSlice, strings.Split(*inputurl, ",")...)
		}
		if *inputfile != "" {
			err, lineSlice := utils.ReadFileAsLine(*inputfile)
			if err != nil {
				utils.OptionsError("Targetfile handle error: "+err.Error(), 2)
			}
			urlSlice = append(urlSlice, lineSlice...)
		}

		if *threads > len(urlSlice) {
			*threads = len(urlSlice)
		}
		if len(urlSlice) == 0 {
			utils.OptionsError("inputurl is empty", 2)
		}
		ctx, cancel := context.WithCancel(context.Background())
		run.Start(ctx, urlSlice, *threads, *timeout, *rulefile)
		//time.Sleep(180 * time.Second)
		cancel()
		pterm.Success.Println("所有任务已经结束。")
	}

	app.Run(os.Args)
}
