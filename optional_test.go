package optional_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/amidgo/optional"
)

func Test_Optional_Empty(t *testing.T) {
	emptyOptional := optional.Optional[int]{}

	value, ok := emptyOptional.Get()

	if value != 0 {
		t.Fatalf("wrong value of empty optional, %d", value)
	}

	if ok {
		t.Fatalf("wrong ok of empty optional, %d", value)
	}

	if !emptyOptional.IsEmpty() {
		t.Fatal("empty optional is not empty")
	}

	if emptyOptional.Pointer() != nil {
		t.Fatal("empty optional pointer not nil")
	}

	driverValue, err := emptyOptional.Value()

	if err != nil {
		t.Fatalf("empty optional Value() err is not nil, %s", err)
	}

	if driverValue != nil {
		t.Fatalf("driver value is not nil, %v", driverValue)
	}
}

func Test_Optional_Comparable(t *testing.T) {
	const initial = 100

	valueOptional := optional.New(initial)

	value, ok := valueOptional.Get()

	if value != initial {
		t.Fatalf("value optional not equal initial, %d", value)
	}

	if !ok {
		t.Fatal("value optional wrong ok")
	}

	if valueOptional.IsEmpty() {
		t.Fatal("value optional is empty")
	}

	pointer := valueOptional.Pointer()

	if pointer == nil {
		t.Fatal("value optional Pointer is nil")
	}

	if *pointer != initial {
		t.Fatal("value optional Pointer not equal initial value")
	}

	driverValue, err := valueOptional.Value()

	if err != nil {
		t.Fatalf("value optional Value err not nil, %s", err)
	}

	i, ok := driverValue.(int)
	if !ok {
		t.Fatalf("invalid type of driverValue, %T", driverValue)
	}

	if i != initial {
		t.Fatalf("value optional is not equal initial, %d", i)
	}

	zeroValueOptional := optional.New(0)

	value, ok = zeroValueOptional.Get()

	if !ok {
		t.Fatal("zero value optional is zero")
	}

	if value != 0 {
		t.Fatalf("zero value optional value wrong value, %d", value)
	}
}

func Test_Optional_Empty_Scan(t *testing.T) {
	scanValue := 100

	emptyOptional := optional.Optional[int]{}

	err := emptyOptional.Scan(scanValue)

	if err != nil {
		t.Fatalf("failed scan to value, %s", err)
	}

	value, ok := emptyOptional.Get()

	if value != scanValue {
		t.Fatalf("scan set wrong value, %d", value)
	}

	if !ok {
		t.Fatalf("scan not affected optional")
	}

	emptyOptional = optional.Optional[int]{}

	err = emptyOptional.Scan("Hello World!")

	if err == nil {
		t.Fatalf("scan from string to int, expected error but got nil")
	}

	value, ok = emptyOptional.Get()

	if ok {
		t.Fatal("scan change emptyOptional")
	}

	if value != 0 {
		t.Fatalf("scan change emptyOptional value, %d", value)
	}
}

func Test_Optional_Value_Scan(t *testing.T) {
	scanValue := 100

	valueOptional := optional.New(1010)

	err := valueOptional.Scan(scanValue)

	if err != nil {
		t.Fatalf("failed scan to value, %s", err)
	}

	value, ok := valueOptional.Get()

	if value != scanValue {
		t.Fatalf("scan set wrong value, %d", value)
	}

	if !ok {
		t.Fatalf("scan not affected optional")
	}

	valueOptional = optional.New(1010)

	err = valueOptional.Scan("Hello World!")

	if err == nil {
		t.Fatalf("scan from string to int, expected error but got nil")
	}

	value, ok = valueOptional.Get()

	if !ok {
		t.Fatal("err scan clear value optional")
	}

	if value != 1010 {
		t.Fatalf("err scan clear previous value from value optional, %d", value)
	}
}

func Test_Optional_MarshalJSON(t *testing.T) {
	emptyOptional := optional.Optional[int]{}

	data, err := emptyOptional.MarshalJSON()

	if err != nil {
		t.Fatalf("empty optional MarshalJSON err not nil, %s", err)
	}

	if !slices.Equal(data, []byte{'n', 'u', 'l', 'l'}) {
		t.Fatalf("empty optional MarshalJSON wrong data, %s", string(data))
	}

	valueOptional := optional.New(100)

	data, err = valueOptional.MarshalJSON()

	if err != nil {
		t.Fatalf("value optional MarshalJSON err not nil, %s", err)
	}

	if !slices.Equal(data, []byte{'1', '0', '0'}) {
		t.Fatalf("value optional MarshalJSON wrong value, %s", string(data))
	}
}

func Test_Optional_UnmarshalJSON(t *testing.T) {
	emptyOptional := optional.Optional[int]{}

	err := emptyOptional.UnmarshalJSON([]byte{'1'})

	if err != nil {
		t.Fatalf("empty optional value UnmarshalJSON err is not nil, %s", err)
	}

	value, ok := emptyOptional.Get()

	if value != 1 {
		t.Fatalf("after UnmarshalJSON receive invalid value, %d", value)
	}

	if !ok {
		t.Fatal("after UnmarshalJSON ok is false")
	}

	valueOptional := optional.New(100)

	err = valueOptional.UnmarshalJSON([]byte{'1'})

	if err != nil {
		t.Fatalf("empty optional value UnmarshalJSON json is not nil, %s", err)
	}

	value, ok = valueOptional.Get()

	if value != 1 {
		t.Fatalf("after UnmarshalJSON receive invalid value, %d", value)
	}

	if !ok {
		t.Fatal("after UnmarshalJSON ok is false")
	}
}

func Test_Optional_Marshal_OmitZero(t *testing.T) {
	omitZeroStruct := struct {
		Value optional.Optional[int] `json:"value,omitzero"`
	}{}

	data, err := json.Marshal(omitZeroStruct)

	if err != nil {
		t.Fatalf("failed marshal omit zero struct, %s", err)
	}

	if !slices.Equal(data, []byte(`{"value":null}`)) {
		t.Fatalf("omit zero marshal invalid data, %s", string(data))
	}
}

func Test_Optional_OmitZero(t *testing.T) {
	valueOptional := optional.New(0).OmitZero()

	value, ok := valueOptional.Get()

	if ok {
		t.Fatal("omit zero optional not omitted")
	}

	if value != 0 {
		t.Fatalf("omit zer optional invalid value, %d", value)
	}
}
