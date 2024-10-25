package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDef(t *testing.T) {
	result, err := EvaluateExpression("(def x (concat \"a\" \"b\"))")
	if err != nil {
		t.Fatalf(err.Error())
	}

	dict := result.(*DictResult)

	got := []string{}
	for _, v := range dict.Values["Variable"] {
		got = append(got, v.String())
	}

	want := []string{"x"}
	assert.Equal(t, want, got)
}
