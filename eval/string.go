package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FurnishString struct{}

func (f *FurnishString) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Furnishing function for string [Internal]"
}

func (f *FurnishString) Symbol() parser.TokLabel {
	return parser.TokString
}

func (f *FurnishString) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() > 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	if stack.Size() == 1 {
		last := stack.Peek()
		if last.Type() != types.TypeCharSequence {
			return WrongTypeOfArg(f.Symbol(), 1, last)
		}
	}

	return nil
}

func (f *FurnishString) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	// pop identifier from stack
	if stack.Size() == 1 {
		last := stack.Pop().(*CharSequenceResult)
		return NewStringResult(last.Val), nil
	}

	return NewStringResult(""), nil
}

func NewFunctionString() (Function, error) {
	fun := &FurnishString{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
