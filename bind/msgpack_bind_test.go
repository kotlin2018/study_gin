package bind

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ugorji/go/codec"
	"testing"
)
func TestMsgpackBindingBindBody(t *testing.T) {
	type teststruct struct {
		Foo string `msgpack:"foo"`
	}
	var s teststruct
	err := msgPackBinding{}.BindBody(msgpackBody(t, teststruct{"FOO"}), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}

func msgpackBody(t *testing.T, obj interface{}) []byte {
	var bs bytes.Buffer
	h := &codec.MsgpackHandle{}
	err := codec.NewEncoder(&bs, h).Encode(obj)
	require.NoError(t, err)
	return bs.Bytes()
}