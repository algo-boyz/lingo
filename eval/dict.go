package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionDict struct{}

func (f *FunctionDict) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(),
			"(p0) ... (pN)",
			parser.TokRightPar),
		"Generates a dictionary from a sequence of pairs (pX) [Internal]"
}

func (f *FunctionDict) Symbol() parser.TokLabel {
	return parser.TokDict
}

func (f *FunctionDict) Validate(env *Environment, stack *StackFrame) error {
	if stack.Empty() {
		return TooFewArgs(f.Symbol(), 0, 1)
	}

	for idx, item := range stack.items {
		if item.Type() != types.TypeDictionary {
			return WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}

	return nil
}

func (f *FunctionDict) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	if stack.Empty() {
		return NewEmptyResult(), fmt.Errorf("(%s) Wrong number of arguments", f.Symbol())
	}

	result := stack.Pop().(*DictResult)
	for !stack.Empty() {
		item := stack.Pop().(*DictResult)
		result = result.Merge(item)
	}

	return result, nil
}

func NewFunctionDict() (Function, error) {
	fun := &FunctionDict{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
