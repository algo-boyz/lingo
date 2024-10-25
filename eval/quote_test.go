package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TestQuoteCommand(t *testing.T) {
	parser.Flush()
	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel(parser.TokString))
	parser.HookToken(parser.TokLabel(parser.TokQuote))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("quote"))
	parser.HookToken(parser.TokLabel("eval"))

	got, err := EvaluateExpression("(quote (concat (concat \"a\" \"b\") \"c\"))")
	assert.NoError(t, err)

	want := `(quote:
 (concat:
  (concat:
   (string "a")
   (string "b")
  )
  (string "c")
 )
)
`
	assert.Equal(t, want, got.String())
	result, err := EvaluateExpression("(eval (quote (concat \"a\" \"b\")))")
	assert.NoError(t, err)
	assert.Equal(t, result.String(), "ab")

	parser.Flush()
}
