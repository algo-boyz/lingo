package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TestRunScript(t *testing.T) {
	parser.Flush()
	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokQuote))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokComment))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel(parser.TokString))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("def"))
	parser.HookToken(parser.TokLabel("resolve"))
	parser.HookToken(parser.TokLabel("bindings"))

	input := `
; this is a lingo script
;; define two variables
(def x "a")
(def y "b")
; this should yield abc
(concat x y "c")`

	result, err := Run(input)
	assert.NoError(t, err)
	assert.Equal(t, result.String(), "abc")
	parser.Flush()
}
