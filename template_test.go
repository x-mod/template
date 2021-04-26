package template

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemplate_Open(t *testing.T) {
	tpl := New(
		Dir("./testdata"),
		NameByPath(true),
		Function("date", ChineseDate),
		Function("rmb", RMB),
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

	name2, err2 := tpl.Find("/404")
	assert.Nil(t, err2)
	assert.Equal(t, name2, "404")

	b := bytes.NewBufferString("")
	err3 := tpl.ExecuteTemplate(b, "a", map[string]interface{}{
		"Now":   time.Now(),
		"Money": 10.58,
	})
	fmt.Println(err3)
	assert.Nil(t, err3)
	fmt.Println(b.String())
}
