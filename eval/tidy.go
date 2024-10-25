package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

type FunctionTidy struct{}

func (f *FunctionTidy) Desc() (string, string) {
	return fmt.Sprintf("%c%s%c",
			parser.TokLeftPar,
			f.Symbol(),
			parser.TokRightPar),
		"Cleanup the environment"
}

func (f *FunctionTidy) Symbol() parser.TokLabel {
	return parser.TokLabel("tidy")
}

func (f *FunctionTidy) Validate(env *Environment, stack *StackFrame) error {
	return nil
}

func (f *FunctionTidy) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	env.Tidy()
	return NewEmptyResult(), nil
}

func NewFunctionTidy() (Function, error) {
	fun := &FunctionTidy{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
