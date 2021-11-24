package run

import (
	"Holmes/utils"
	"fmt"
)

type Responseinfo struct {
	Url				string
	Content     	string
	Headers			string
	Server			string
	Title			string
	Cert			string
}

func (e *Engine) Scanner(taskurlInterface interface{}) {
	defer e.Wg.Done()
	taskurl := taskurlInterface.(string)
	responinfo, err := utils.Requrl(taskurl, e.Options.Timeout)
	if err != nil {
		fmt.Printf("can't request target: " + taskurl + "\n")
	}
	if products, err := e.ExecuteCel(Responseinfo(responinfo), rules_global); err == nil {
		if err == nil {
			fmt.Println(taskurl + " : " + products)
		}
	}
}
