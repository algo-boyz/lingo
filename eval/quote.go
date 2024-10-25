package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionQuote struct{}

func (f *FunctionQuote) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "x",
			parser.TokRightPar),
		"quote expression s"
}

func (f *FunctionQuote) Symbol() parser.TokLabel {
	return parser.TokLabel("quote")
}

func (f *FunctionQuote) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	first := stack.GetArgument(0)
	if first.Type().Id != types.TypeSexpressionId {
		return WrongTypeOfArg(f.Symbol(), 1, first)
	}

	return nil
}

func (f *FunctionQuote) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	sexp := stack.Pop()
	return sexp, nil
}

func NewFunctionQuote() (Function, error) {
	fun := &FunctionQuote{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
