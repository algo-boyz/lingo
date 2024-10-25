package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

// this function is useful to extract comments
type FunctionComment struct{}

func (f *FunctionComment) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"help with the evaluation of comnments [Internal]"
}

func (f *FunctionComment) Symbol() parser.TokLabel {
	return parser.TokLabel(parser.TokComment)
}

func (f *FunctionComment) Validate(env *Environment, stack *StackFrame) error {
	return nil
}

func (f *FunctionComment) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	return stack.Pop(), nil
}

func NewFunctionComment() (Function, error) {
	fun := &FunctionComment{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
