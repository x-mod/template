package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate_Open(t *testing.T) {
	tpl := New(
		Dir("./testdata"),
		NameByPath(true),
	)
	assert.Nil(t, tpl.Open())
	assert.NotNil(t, tpl.Lookup("a"))

}
