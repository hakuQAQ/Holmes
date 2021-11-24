package run

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
)

func Start(ctx context.Context, inputurl []string,threads int, timeout int, rulefile string) {
	//退出前保存结果
	//go Finderlist()

	option := NewDefaultOptions(inputurl, threads, timeout, rulefile)
	engine := CreateEngine(option)

	if err := engine.Parser(); err != nil {
		fmt.Println(err)
	}

	engine.Detect()
	select {
	case <-ctx.Done():
		pterm.Success.Println("父线程已经退去")
	default:
		engine.Wg.Wait()
	}
}