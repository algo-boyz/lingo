package main

import (
	"fmt"
	"os"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/eval"
)

func main() {
	if len(os.Args) <= 1 {
		eval.RunLoop()
		os.Exit(0)
	}

	result, err := eval.RunScriptPath(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(result)

	os.Exit(0)
}
