package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

func AddTest(t *testing.T) {
	result, err := EvaluateExpression("(add (add 1 2 3 4) (add 1 2 3 4) 5 6)")
	assert.NoError(t, err)
	if result.Type().HasProperty(types.Primitive) {
		t.Errorf("Wrong type")
	}
	intResult := result.(*IntResult)
	assert.Equal(t, intResult.Val, 31)
}
