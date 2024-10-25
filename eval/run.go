package eval

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const Prompt = "\033[32m\033[0m > "

func RunLoop() {
	out := os.Stdout
	in := os.Stdin

	scanner := bufio.NewScanner(in)
	env := NewEnvironment()

	defer env.GC()

	for {
		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			io.WriteString(out, "Parsing error")
		}

		line := scanner.Text()

		if line == "exit" || line == "quit" {
			io.WriteString(out, "Bye")
			io.WriteString(out, "\n")
			break
		}

		result, err := EvaluateExpressionWithEnv(line, env)
		if err != nil {
			io.WriteString(out, err.Error())
			io.WriteString(out, "\n")
			continue
		}
		io.WriteString(out, result.String())
		io.WriteString(out, "\n")
	}
}

func Run(scriptcontent string) (Result, error) {
	result, err := EvaluateExpression(scriptcontent)
	if err != nil {
		return nil, err
	}

	return result, err
}

func RunScriptPath(scriptpath string) (Result, error) {
	content, err := os.ReadFile(scriptpath)
	if err != nil {
		return nil, err
	}
	return Run(string(content))
}
