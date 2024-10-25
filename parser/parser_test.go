package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSingleLine(t *testing.T) {
	Flush()
	HookToken(TokLabel(TokLeftPar))
	HookToken(TokLabel(TokRightPar))
	HookToken(TokLabel(TokQuote))
	HookToken(TokLabel(TokIdentifier))
	HookToken(TokLabel(TokEof))
	HookToken(TokLabel("concat"))
	HookToken(TokLabel("string"))

	input := "(concat (concat \"a\") \"b\")"

	parser := NewParser()
	got, err := parser.Parse(input)
	assert.NoError(t, err)

	want := `(root:
 (concat:
  (concat:
   (string "a")
  )
  (string "b")
 )
)
`
	assert.Equal(t, want, got.String())
	Flush()
}

func TestParseMultiLine(t *testing.T) {
	Flush()
	HookToken(TokLabel(TokLeftPar))
	HookToken(TokLabel(TokRightPar))
	HookToken(TokLabel(TokQuote))
	HookToken(TokLabel(TokIdentifier))
	HookToken(TokLabel(TokEof))
	HookToken(TokLabel("concat"))
	HookToken(TokLabel("string"))

	input := `(concat (concat "x") "y")
(concat (concat "a") "b")`

	parser := NewParser()
	got, err := parser.Parse(input)
	assert.NoError(t, err)

	want := `(root:
 (concat:
  (concat:
   (string "x")
  )
  (string "y")
 )
 (concat:
  (concat:
   (string "a")
  )
  (string "b")
 )
)
`

	assert.Equal(t, want, got.String())
	Flush()
}
