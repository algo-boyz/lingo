package eval

import (
	"fmt"
	"strconv"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FurnishInt struct{}

func (f *FurnishInt) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Furnishing function for string [Internal]"
}

func (f *FurnishInt) Symbol() parser.TokLabel {
	return parser.TokInt
}

func (f *FurnishInt) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	last := stack.Peek()
	if last.Type() != types.TypeCharSequence {
		return WrongTypeOfArg(f.Symbol(), 1, last)
	}

	value := last.(*CharSequenceResult)
	_, err := strconv.Atoi(value.Val)
	if err != nil {
		return err
	}

	return nil
}

func (f *FurnishInt) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	// pop identifier from stack
	last := stack.Pop().(*CharSequenceResult)
	value, _ := strconv.Atoi(last.Val)

	return NewIntResult(value), nil
}

func NewFunctionInt() (Function, error) {
	fun := &FurnishInt{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
