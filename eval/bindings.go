package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

type FunctionBindings struct{}

func (f *FunctionBindings) Desc() (string, string) {
	return fmt.Sprintf("%c%s%c",
			parser.TokLeftPar,
			f.Symbol(),
			parser.TokRightPar),
		"List all bindings in the current environment"
}

func (f *FunctionBindings) Symbol() parser.TokLabel {
	return parser.TokLabel("bindings")
}

func (f *FunctionBindings) Validate(env *Environment, stack *StackFrame) error {
	return nil
}

func (f *FunctionBindings) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	variable := "Variable"
	typ := "Type"

	result := NewDictResult(variable, typ)

	for symbol, item := range env.Bindings {
		result.AddPair(typ, NewStringResult(item.Type().Name))
		result.AddPair(variable, NewStringResult(symbol))
	}

	return result, nil
}

func NewFunctionBindings() (Function, error) {
	fun := &FunctionBindings{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
