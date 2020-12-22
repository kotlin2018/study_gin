package bind

import (
	"bytes"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testInterface interface {
	String() string
}

type subStructNoValidation struct {
	IString string
	IInt    int
}

type mapNoValidationSub map[string]subStructNoValidation

type Object map[string]interface{}

// 结构体未验证的值
type structNoValidationValues struct {
	subStructNoValidation
	Boolean bool
	UInteger   uint
	Integer    int
	Integer8   int8
	Integer16  int16
	Integer32  int32
	Integer64  int64
	Uinteger8  uint8
	Uinteger16 uint16
	Uinteger32 uint32
	Uinteger64 uint64

	Float32 float32
	Float64 float64

	String string

	Date time.Time

	Struct        subStructNoValidation
	InlinedStruct struct {
		String  []string
		Integer int
	}

	IntSlice           []int
	IntPointerSlice    []*int
	StructPointerSlice []*subStructNoValidation
	StructSlice        []subStructNoValidation
	InterfaceSlice     []testInterface

	UniversalInterface interface{}
	CustomInterface    testInterface

	FloatMap  map[string]float32
	StructMap mapNoValidationSub
}

func createNoValidationValues() structNoValidationValues {
	integer := 1
	s := structNoValidationValues{
		Boolean:            true,
		UInteger:           1 << 29,
		Integer:            -10000,
		Integer8:           120,
		Integer16:          -20000,
		Integer32:          1 << 29,
		Integer64:          1 << 61,
		Uinteger8:          250,
		Uinteger16:         50000,
		Uinteger32:         1 << 31,
		Uinteger64:         1 << 62,
		Float32:            123.456,
		Float64:            123.456789,
		String:             "text",
		Date:               time.Time{},
		CustomInterface:    &bytes.Buffer{},
		Struct:             subStructNoValidation{},
		IntSlice:           []int{-3, -2, 1, 0, 1, 2, 3},
		IntPointerSlice:    []*int{&integer},
		StructSlice:        []subStructNoValidation{},
		UniversalInterface: 1.2,
		FloatMap: map[string]float32{
			"foo": 1.23,
			"bar": 232.323,
		},
		StructMap: mapNoValidationSub{
			"foo": subStructNoValidation{},
			"bar": subStructNoValidation{},
		},
		// StructPointerSlice []noValidationSub
		// InterfaceSlice     []testInterface
	}
	s.InlinedStruct.Integer = 1000
	s.InlinedStruct.String = []string{"first", "second"}
	s.IString = "substring"
	s.IInt = 987654
	return s
}

func TestValidateNoValidationValues(t *testing.T) {
	origin := createNoValidationValues()
	test := createNoValidationValues()
	empty := structNoValidationValues{}

	assert.Nil(t, validate(test))
	assert.Nil(t, validate(&test))
	assert.Nil(t, validate(empty))
	assert.Nil(t, validate(&empty))

	assert.Equal(t, origin, test)
}

func TestValidatePrimitives(t *testing.T) {
	obj := Object{"foo": "bar", "bar": 1}
	assert.NoError(t, validate(obj))
	assert.NoError(t, validate(&obj))
	assert.Equal(t, Object{"foo": "bar", "bar": 1}, obj)

	obj2 := []Object{{"foo": "bar", "bar": 1}, {"foo": "bar", "bar": 1}}
	assert.NoError(t, validate(obj2))
	assert.NoError(t, validate(&obj2))

	nu := 10
	assert.NoError(t, validate(nu))
	assert.NoError(t, validate(&nu))
	assert.Equal(t, 10, nu)

	str := "value"
	assert.NoError(t, validate(str))
	assert.NoError(t, validate(&str))
	assert.Equal(t, "value", str)
}

// structCustomValidation是我们用来检查的助手结构体
// 可以在其上注册自定义验证。
// notOne 是自定义验证器
type structCustomValidation struct {
	Integer int `validator:"notone"`
}

func notOne(f1 validator.FieldLevel) bool {
	if val, ok := f1.Field().Interface().(int); ok {
		return val != 1
	}
	return false
}

func TestValidatorEngine(t *testing.T) {
	//验证函数' notOne '是否匹配
	// defaultValidator指定的函数签名
	//通过扩展验证器库。
	engine, ok := Validator.Engine().(*validator.Validate)
	assert.True(t, ok)

	err := engine.RegisterValidation("notone", notOne)
	// Check that we can register custom validation without error
	assert.Nil(t, err)

	// Create an instance which will fail validation
	withOne := structCustomValidation{Integer: 1}
	errs := validate(withOne)

	// Check that we got back non-nil errs
	assert.NotNil(t, errs)
	// Check that the error matches expectation
	assert.Error(t, errs, "", "", "notOne")
}
