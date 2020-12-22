package bind

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonBinding_BindBody(t *testing.T) {
	var s struct{
		Foo string   `json:"foo"`
	}
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO"}`), &s)
	require.NoError(t,err)
	assert.Equal(t,"FOO",s.Foo)

	//jsonBinding{}.BindBody([]byte(`{"foo": "1111"}`), &s)
	//fmt.Println(s.Foo)
}

func TestJSONBindingBindBodyMap(t *testing.T) {
	s := make(map[string]string)
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO","hello":"world"}`), &s)
	require.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Equal(t, "FOO", s["foo"])
	assert.Equal(t, "world", s["hello"])

	//jsonBinding{}.BindBody([]byte(`{"foo": "9999","hello":"world"}`), &s)
	//require.NoError(t, err)
	//assert.Len(t, s, 2)
	//assert.Equal(t, "FOO", s["foo"])
	//assert.Equal(t, "world", s["hello"])
	//fmt.Println(s["foo"])
	//fmt.Println(s["hello"])
}
