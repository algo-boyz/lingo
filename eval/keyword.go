package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FurnishKeyword struct{}

func (f *FurnishKeyword) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Furnishing function for keywords [Internal]"
}

func (f *FurnishKeyword) Symbol() parser.TokLabel {
	return parser.TokKeyword
}

func (f *FurnishKeyword) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	last := stack.Peek()
	if last.Type() != types.TypeCharSequence {
		return WrongTypeOfArg(f.Symbol(), 1, last)
	}

	return nil
}

func (f *FurnishKeyword) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	// pop identifier from stack
	last := stack.Pop().(*CharSequenceResult)
	return NewKeyWordResult(last.Val), nil
}

func NewFunctionKeyword() (Function, error) {
	fun := &FurnishKeyword{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
