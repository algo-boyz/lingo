package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionDef struct{}

func (f *FunctionDef) Desc() (string, string) {
	return fmt.Sprintf("%c%s%s%c",
			parser.TokLeftPar,
			f.Symbol(), " a b",
			parser.TokRightPar),
		"Declare variable a and bind it to the result yielded by expression b"
}

func (f *FunctionDef) Symbol() parser.TokLabel {
	return parser.TokLabel("def")
}

func (f *FunctionDef) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() < 2 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 2)
	}
	first := stack.GetArgument(0)
	if first.Type() != types.TypeSymbol {
		return WrongTypeOfArg(f.Symbol(), 1, first)
	}

	err := CheckVariableName(f.Symbol(), 1, first)
	if err != nil {
		return nil
	}

	return nil
}

func (f *FunctionDef) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	second := stack.Pop()
	first := stack.Pop().(*SymbolResult)

	env.Bind(first.String(), second)

	variable := "Variable"
	typ := "Type"
	value := "Value"

	result := NewDictResult(variable, typ, value)
	if err := result.AddPair(variable, first); err != nil {
		return nil, err
	}
	if err := result.AddPair(value, second); err != nil {
		return nil, err
	}
	if err := result.AddPair(typ, NewStringResult(second.Type().Name)); err != nil {
		return nil, err
	}

	return result, nil
}

func NewFunctionDef() (Function, error) {
	fun := &FunctionDef{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
