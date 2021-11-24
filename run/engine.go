package run

import (
	"Holmes/utils"
	"github.com/panjf2000/ants"
	"github.com/pterm/pterm"
	"sync"
)

type Options struct {
	Url			[]string
	Threads		int
	Timeout		int
	Rulefile		string
}

type Engine struct {
	TaskUrl			[]string
	ScanInfo		map[string]string
	Wg              *sync.WaitGroup
	Options         *Options
}

var rules_global []YamlRule

func NewDefaultOptions(inputurl []string,threads int,timeout int, rulefile string) *Options {
	return &Options{
		Url:		inputurl,
		Threads: 	threads,
		Timeout: 	timeout,
		Rulefile: 	rulefile,
	}
}

//创建引擎
func CreateEngine(option *Options) *Engine {
	return &Engine{
		Wg:          &sync.WaitGroup{},
		Options:     option,
	}
}

//todo:将ip端口范围转换为切片
func (e *Engine) Parser() error {
	for _, url := range e.Options.Url {
		e.TaskUrl = append(e.TaskUrl, url)
	}
	return nil
}

func (e *Engine) Detect() {
	pool, err := ants.NewPoolWithFunc(e.Options.Threads, e.Scanner)
	if err != nil {
		utils.OptionsError("NewPool err", 2)
	}
	pterm.Info.Println("正在读取yaml文件... (请等待)")
	rules_global, err = LoadYaml(e.Options.Rulefile)
	if err != nil {
		utils.OptionsError("LoadYaml err", 2)
	}
	defer pool.Release()
	pterm.Info.Println("正在进行指纹识别... (请等待)")
	for _, url := range e.TaskUrl {
		e.Wg.Add(1)
		pool.Invoke(url)
	}
}



