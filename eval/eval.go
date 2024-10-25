package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/traversal"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionEval struct{}

func (f *FunctionEval) Desc() (string, string) {
	return fmt.Sprintf("%c%s%s%c",
			parser.TokLeftPar,
			f.Symbol(), " s",
			parser.TokRightPar),
		"evaluate sub-expression s"
}

func (f *FunctionEval) Symbol() parser.TokLabel {
	return parser.TokLabel("eval")
}

func (f *FunctionEval) Validate(env *Environment, stack *StackFrame) error {
	if stack.Size() != 1 {
		return WrongNumberOfArgs(f.Symbol(), stack.Size(), 1)
	}
	first := stack.GetArgument(0)
	if first.Type() != types.TypeSexpression {
		return WrongTypeOfArg(f.Symbol(), 1, first)
	}

	return nil
}

func (f *FunctionEval) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	expression := stack.Pop().(*SexpressionResult)
	expressionToEvaluate := expression.Exp

	if expression.Exp.Kind == funquote.Symbol() {
		// unwrap, unquote
		expressionToEvaluate = expression.Exp.SubExpressions[0]
	}

	evaluator := NewEvaluator(env)
	visitor := traversal.NewExpressionWalker(&evaluator)
	err := visitor.Walk(expressionToEvaluate)
	if err != nil {
		return nil, err
	}
	return evaluator.Result(), nil
}

func NewFunctionEval() (Function, error) {
	fun := &FunctionEval{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
