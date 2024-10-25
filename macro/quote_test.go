package macro

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TestQuote(t *testing.T) {
	parser.Flush()

	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel(parser.TokSingleQuoteChar))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("string"))
	parser.HookToken(parser.TokLabel("+"))

	input := "'(+ 13 4 4 5)"

	p := parser.NewParser()
	parsed, err := p.Parse(input)
	assert.NoError(t, err)

	got, err := ApplyMacros(parsed)
	assert.NoError(t, err)

	want := `(root:
 (quote:
  (+:
   (int 13)
   (int 4)
   (int 4)
   (int 5)
  )
 )
)
`
	assert.Equal(t, want, got.String())
	parser.Flush()
}
