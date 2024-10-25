package eval

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

func TestTidy(t *testing.T) {
	parser.Flush()
	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokQuote))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel(parser.TokString))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("def"))
	parser.HookToken(parser.TokLabel("resolve"))
	parser.HookToken(parser.TokLabel("bindings"))
	parser.HookToken(parser.TokLabel("tidy"))

	env := NewEnvironment()
	_, err := EvaluateExpressionWithEnv("(def x (concat \"a\" \"b\"))", env)
	assert.NoError(t, err)

	_, err = EvaluateExpressionWithEnv("(def y 1)", env)
	assert.NoError(t, err)

	result, err := EvaluateExpressionWithEnv("(bindings)", env)
	assert.NoError(t, err)

	assert.Equal(t, types.TypeDictionary, result.Type())
	dict := result.(*DictResult)

	got := []string{}
	for _, v := range dict.Values["Variable"] {
		got = append(got, v.String())
	}

	want := []string{"x", "y"}

	sort.Strings(got)
	sort.Strings(want)
	assert.Equal(t, want, got)

	_, err = EvaluateExpressionWithEnv("(tidy)", env)
	assert.NoError(t, err)

	result, err = EvaluateExpressionWithEnv("(bindings)", env)
	assert.NoError(t, err)

	assert.Equal(t, types.TypeDictionary, result.Type())
	dict = result.(*DictResult)

	assert.Len(t, dict.Values["Variable"], 0)

	parser.Flush()
}
