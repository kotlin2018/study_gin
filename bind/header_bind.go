package bind

import (
	"net/http"
	"net/textproto"
	"reflect"
)

type headerBinding struct {}
type headerSource map[string][]string
var _ setter = headerSource(nil)

func (hs headerSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSet bool, err error) {
	return setByForm(value, field, hs, textproto.CanonicalMIMEHeaderKey(tagValue), opt)
}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(r *http.Request, obj interface{}) error {
	if err := mapHeader(obj, r.Header); err != nil {
		return err
	}
	return validate(obj)
}

func mapHeader(ptr interface{}, h map[string][]string) error {
	return mappingByPtr(ptr, headerSource(h), "header")
}