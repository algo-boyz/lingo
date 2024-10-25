package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FurnishIdentifier struct{}

func (f *FurnishIdentifier) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Returns an identifier handle [Internal]"
}

func (f *FurnishIdentifier) Symbol() parser.TokLabel {
	return parser.TokIdentifier
}

func (f *FurnishIdentifier) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	last := stack.Peek()
	if last.Type() != types.TypeCharSequence {
		return WrongTypeOfArg(f.Symbol(), 1, last)
	}

	return nil
}

func (f *FurnishIdentifier) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	// pop identifier from stack
	last := stack.Pop().(*CharSequenceResult)
	return NewSymbolResult(last.Val), nil
}

func NewFunctionIdentifier() (Function, error) {
	fun := &FurnishIdentifier{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
