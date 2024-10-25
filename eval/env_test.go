package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDict(t *testing.T) {
	result := NewDictResult("name", "profession")
	result.Values["name"] = append(result.Values["name"], NewStringResult("Anna"), NewStringResult("John"))
	result.Values["profession"] = append(result.Values["profession"], NewStringResult("Writer"), NewStringResult("Driver"))

	expected := `+------+------------+
| NAME | PROFESSION |
+------+------------+
| Anna | Writer     |
+------+------------+
| John | Driver     |
+------+------------+
`

	assert.Equal(t, expected, result.String())
}
