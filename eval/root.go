package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

type FunctionRoot struct{}

func (f *FunctionRoot) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Evaluates the subexpression and propagates the result [Internal]"
}

func (f *FunctionRoot) Symbol() parser.TokLabel {
	return parser.TokRoot
}

func (f *FunctionRoot) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() == 0 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}
	return nil
}

func (f *FunctionRoot) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	return stack.Pop(), nil
}

func NewFunctionRoot() (Function, error) {
	fun := &FunctionRoot{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
