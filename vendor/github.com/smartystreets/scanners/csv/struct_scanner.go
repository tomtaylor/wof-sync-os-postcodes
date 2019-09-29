package csv

import (
	"fmt"
	"io"
	"reflect"
)

type StructScanner struct {
	*ColumnScanner
}

func NewStructScanner(reader io.Reader, options ...Option) (*StructScanner, error) {
	inner, err := NewColumnScanner(reader, options...)
	if err != nil {
		return nil, err
	}
	return &StructScanner{ColumnScanner: inner}, nil
}

func (this *StructScanner) Populate(v interface{}) error {
	type_ := reflect.TypeOf(v)
	if type_.Kind() != reflect.Ptr {
		return fmt.Errorf("Provided value must be reflect.Ptr. You provided [%v] ([%v]).", v, type_.Kind())
	}

	value := reflect.ValueOf(v)
	if value.IsNil() {
		return fmt.Errorf("The provided value was nil. Please provide a non-nil pointer.")
	}

	this.populate(type_.Elem(), value.Elem())
	return nil
}

func (this *StructScanner) populate(type_ reflect.Type, value reflect.Value) {
	for x := 0; x < type_.NumField(); x++ {
		column := type_.Field(x).Tag.Get("csv")

		_, found := this.columnIndex[column]
		if !found {
			continue
		}

		field := value.Field(x)
		if field.Kind() != reflect.String {
			continue // Future: return err?
		} else if !field.CanSet() {
			continue // Future: return err?
		}

		field.SetString(this.Column(column))
	}
}
