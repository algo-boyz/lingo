package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionResolve struct{}

func (f *FunctionResolve) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(), "a",
			parser.TokRightPar),
		"Resolves variable a [Internal]"
}

func (f *FunctionResolve) Symbol() parser.TokLabel {
	return parser.TokLabel("resolve")
}

func (f *FunctionResolve) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}

	last := stack.Peek()
	if last.Type() != types.TypeCharSequence {
		return WrongTypeOfArg(f.Symbol(), 1, last)
	}

	if _, ok := env.Bindings[last.String()]; !ok {
		return fmt.Errorf(ErrorMessage(f.Symbol(), "variable '%s' not defined"), last.String())
	}

	return nil
}

func (f *FunctionResolve) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	last := stack.Pop()
	if v, ok := env.Bindings[last.String()]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("(%s) variable '%s' not defined", f.Symbol(), last.String())
}

func NewFunctionResolve() (Function, error) {
	fun := &FunctionResolve{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
