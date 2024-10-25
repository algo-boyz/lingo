package eval

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

func TestBindings(t *testing.T) {
	env := NewEnvironment()
	_, err := EvaluateExpressionWithEnv("(def x (concat \"a\" \"b\"))", env)
	assert.NoError(t, err)

	_, err = EvaluateExpressionWithEnv("(def y 1)", env)
	assert.NoError(t, err)

	result, err := EvaluateExpressionWithEnv("(bindings)", env)
	assert.NoError(t, err)
	assert.Equal(t, result.Type(), types.TypeDictionary)
	dict := result.(*DictResult)

	got := []string{}
	for _, v := range dict.Values["Variable"] {
		got = append(got, v.String())
	}

	want := []string{"x", "y"}
	sort.Strings(got)
	sort.Strings(want)
	assert.Equal(t, want, got)
}
