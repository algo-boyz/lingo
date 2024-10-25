package eval

import (
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/macro"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/traversal"
)

func EvaluateExpression(expression string) (Result, error) {
	return EvaluateExpressionWithEnv(expression, NewEnvironment())
}

func EvaluateExpressionWithEnv(expression string, env *Environment) (Result, error) {
	parser := parser.NewParser()

	p, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}

	p, err = macro.ApplyMacros(p)
	if err != nil {
		return nil, err
	}

	evaluator := NewEvaluator(env)
	visitor := traversal.NewExpressionWalker(&evaluator)
	err = visitor.Walk(p)
	if err != nil {
		return nil, err
	}

	return evaluator.Result(), nil
}
