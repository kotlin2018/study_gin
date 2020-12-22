package bind

import (
	"fmt"
	"testing"
)

func TestXMLBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `xml:"foo"`
	}
	xmlBody := `<?xml version="1.0" encoding="UTF-8"?>
<root>
   <foo>1111</foo>
</root>`
	//err := xmlBinding{}.BindBody([]byte(xmlBody), &s)
	//require.NoError(t, err)
	//assert.Equal(t, "FOO", s.Foo)

	xmlBinding{}.BindBody([]byte(xmlBody), &s)
	fmt.Println(s.Foo)
	// 输出 1111
}