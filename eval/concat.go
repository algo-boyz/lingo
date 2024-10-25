package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionConcat struct{}

func (f *FunctionConcat) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a b",
			parser.TokRightPar),
		"Concatenates string sub-expressions a and b"
}

func (f *FunctionConcat) Symbol() parser.TokLabel {
	return parser.TokLabel("concat")
}

func (f *FunctionConcat) Validate(env *Environment, stack *StackFrame) error {
	if stack.Empty() {
		return TooFewArgs(f.Symbol(), 0, 1)
	}
	for idx, item := range stack.items {
		if item.Type() != types.TypeString {
			return WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}
	return nil
}

func (f *FunctionConcat) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	result := NewStringResult("")
	for !stack.Empty() {
		item := stack.Pop()
		result.Val = item.String() + result.Val
	}

	return result, nil
}

func NewFunctionConcat() (Function, error) {
	fun := &FunctionConcat{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
