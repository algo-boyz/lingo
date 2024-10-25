package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionDesc struct{}

func (f *FunctionDesc) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(),
			":x",
			parser.TokRightPar),
		"Provide a description for the function identified by the keyword x"
}

func (f *FunctionDesc) Symbol() parser.TokLabel {
	return parser.TokLabel("desc")
}

func (f *FunctionDesc) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() > 1 {
		return TooManyArgs(f.Symbol(), stack.Size(), 1)
	}

	if stack.Size() == 1 && stack.Peek().Type() != types.TypeKeyword {
		return WrongTypeOfArg(f.Symbol(), 1, stack.Peek())
	}

	return nil
}

func (f *FunctionDesc) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	cmdKey := "Command"
	descKey := "Description"
	result := NewDictResult(cmdKey, descKey)

	if stack.Size() == 0 {
		for _, builtin := range builtins {
			cmd, desc := builtin.Desc()
			result.AddPair(cmdKey, NewStringResult(cmd))
			result.AddPair(descKey, NewStringResult(desc))
		}
	} else {
		last := stack.Pop()
		builtin, ok := builtins[parser.TokLabel(last.String())]
		if !ok {
			return nil, fmt.Errorf(ErrorMessage(f.Symbol(),
				"could not find functino '%s'"), last.String())
		}
		cmd, desc := builtin.Desc()
		result.AddPair(cmdKey, NewStringResult(cmd))
		result.AddPair(descKey, NewStringResult(desc))
	}

	return result, nil
}

func NewFunctionDesc() (Function, error) {
	fun := &FunctionDesc{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
