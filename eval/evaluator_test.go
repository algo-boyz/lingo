package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/traversal"
)

func TestSexpressionWithOption(t *testing.T) {
	parser.Flush()
	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokQuote))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokLeftCurly))
	parser.HookToken(parser.TokLabel(parser.TokRightCurly))
	parser.HookToken(parser.TokLabel(parser.TokKeyword))
	parser.HookToken(parser.TokLabel(parser.TokLeftBracket))
	parser.HookToken(parser.TokLabel(parser.TokRightBracket))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel(parser.TokDict))
	parser.HookToken(parser.TokLabel(parser.TokString))
	parser.HookToken(parser.TokLabel(parser.TokComma))

	// support for support arrays in dictionaries
	input := "{ :optiona [ 1 2 ], :optionb [ \"a\" :b ] }"

	p := parser.NewParser()
	got, err := p.Parse(input)
	assert.NoError(t, err)

	env := NewEnvironment()
	evaluator := NewEvaluator(env)
	visitor := traversal.NewExpressionWalker(&evaluator)
	err = visitor.Walk(got)
	assert.NoError(t, err)

	want := `+----------+----------+
| :OPTIONB | :OPTIONA |
+----------+----------+
| a,:b     | 1,2      |
+----------+----------+
`

	assert.Equal(t, want, evaluator.Result().String())
	parser.Flush()
}
