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

	name, err := tpl.Find("/sub/c/sdfs")
	assert.Nil(t, err)
	assert.Equal(t, name, "sub/c")

	_, err = tpl.Find("/e/c/sdfs")
	assert.NotNil(t, err)

	name1, err1 := tpl.Find("/")
	assert.Nil(t, err1)
	assert.Equal(t, name1, "index")

}
