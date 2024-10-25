package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

func ConcatTest(t *testing.T) {
	result, err := EvaluateExpression("(concat (concat (concat \"a\") \"b\") \"c\")")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result.Type() != types.TypeString {
		t.Errorf("Wrong type")
	}
	strResult := result.(*StringResult)
	assert.Equal(t, "abc", strResult.Val)
}
