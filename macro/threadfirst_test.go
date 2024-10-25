package macro

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TestThreadingFirstMath(t *testing.T) {
	parser.Flush()

	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("string"))
	parser.HookToken(parser.TokLabel("->"))
	parser.HookToken(parser.TokLabel("-"))
	parser.HookToken(parser.TokLabel("+"))
	parser.HookToken(parser.TokLabel("*"))

	// This is thread first macro which improves readability because it turns
	// the Sexpression inside out and, thus reduces the "nestedness"
	// The semantics is similar to a unix pipe
	// where the result of the first expression is piped into the next one
	// and treated as the first parameter
	// (-> (+ 13 4) (- 5) (* 3)) ==> (- 17 5) (* 3) ==> (* 12 3) ==> 36
	// The expression below is syntactic sugar for (* (- (+ 13 4) 5) 3)
	input := "(-> (+ 13 4) (- 5) (* 3))"

	p := parser.NewParser()
	parsed, err := p.Parse(input)
	assert.NoError(t, err)

	got, err := ApplyMacros(parsed)
	assert.NoError(t, err)

	want := `(root:
 (*:
  (-:
   (+:
    (int 13)
    (int 4)
   )
   (int 5)
  )
  (int 3)
 )
)
`
	assert.Equal(t, want, got.String())
	parser.Flush()
}

func TestThreadingFirstString(t *testing.T) {
	parser.Flush()

	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
	parser.HookToken(parser.TokLabel(parser.TokIdentifier))
	parser.HookToken(parser.TokLabel(parser.TokEof))
	parser.HookToken(parser.TokLabel(parser.TokInt))
	parser.HookToken(parser.TokLabel("concat"))
	parser.HookToken(parser.TokLabel("string"))
	parser.HookToken(parser.TokLabel("->"))

	input := "(-> \"hello\" (concat \"this\") (concat \"is\") (concat \"a\") (concat \"test\"))"

	p := parser.NewParser()
	parsed, err := p.Parse(input)
	assert.NoError(t, err)

	got, err := ApplyMacros(parsed)
	assert.NoError(t, err)

	want := `(root:
 (concat:
  (concat:
   (concat:
    (concat:
     (string "hello")
     (string "this")
    )
    (string "is")
   )
   (string "a")
  )
  (string "test")
 )
)
`
	assert.Equal(t, want, got.String())
	parser.Flush()
}
