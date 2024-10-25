package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSexpression(t *testing.T) {

	Flush()
	HookToken(TokLabel(TokLeftPar))
	HookToken(TokLabel(TokRightPar))
	HookToken(TokLabel(TokQuote))
	HookToken(TokLabel(TokIdentifier))
	HookToken(TokLabel(TokEof))
	HookToken(TokLabel("concat"))
	HookToken(TokLabel("string"))

	input := "(concat (concat \"\") \"b\")"

	parser := NewParser()
	got, err := parser.Parse(input)
	assert.NoError(t, err)

	want := `(root:
 (concat:
  (concat:
   (string "")
  )
  (string "b")
 )
)
`

	assert.Equal(t, want, got.String())
	assert.Equal(t, want, got.DeepCopy().String())

	Flush()
}
