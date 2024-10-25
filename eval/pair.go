package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionPair struct{}

func (f *FunctionPair) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), ":a b",
			parser.TokRightPar),
		"Generate a data pair with the symbol a and the result of the (string|int) sub-expression b [Internal]"
}

func (f *FunctionPair) Symbol() parser.TokLabel {
	return parser.TokLabel(parser.TokPair)
}

func (f *FunctionPair) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 2 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 2)
	}

	first := stack.GetArgument(0)
	if first.Type().Id != types.TypeKeywordId {
		return WrongTypeOfArg(f.Symbol(), 1, first)
	}

	return nil
}

func (f *FunctionPair) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	second := stack.Pop()
	first := stack.Pop().(*KeyWordResult)

	result := NewDictResult(first.String())
	if err := result.AddPair(first.String(), second); err != nil {
		return nil, err
	}
	return result, nil
}

func NewFunctionPair() (Function, error) {
	fun := &FunctionPair{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
