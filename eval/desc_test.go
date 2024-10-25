package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

func TestDesc(t *testing.T) {
	result, err := EvaluateExpression("(desc)")
	assert.NoError(t, err)
	assert.Equal(t, types.TypeDictionary, result.Type())
	if result.Type() != types.TypeDictionary {
		t.Errorf("Wrong type")
	}
	dict := result.(*DictResult)
	assert.Len(t, dict.Values["Command"], len(builtins))
}
