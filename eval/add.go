package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionAdd struct{}

func (f *FunctionAdd) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(),
			"a b",
			parser.TokRightPar),
		"Add numeric sub-expressions a and b"
}

func (f *FunctionAdd) Symbol() parser.TokLabel {
	return parser.TokLabel("add")
}

func (f *FunctionAdd) Validate(env *Environment, stack *StackFrame) error {
	if stack.Empty() {
		return TooFewArgs(f.Symbol(), 0, 1)
	}
	for idx, item := range stack.items {
		if !item.Type().HasProperty(types.Numeric) {
			return WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}
	return nil
}

func (f *FunctionAdd) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	result := NewIntResult(0)
	for !stack.Empty() {
		item := stack.Pop().(*IntResult)
		result.Val = item.Val + result.Val
	}

	return result, nil
}

func NewFunctionAdd() (Function, error) {
	fun := &FunctionAdd{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
