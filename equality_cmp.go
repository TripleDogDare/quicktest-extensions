package qtx

import (
	"errors"
	"fmt"
	"reflect"

	qt "github.com/frankban/quicktest"
)

var LessThan = equalityCmp{lessThan: true}
var LessThanOrEqual = equalityCmp{lessThan: true, allowEqual: true}
var GreaterThan = equalityCmp{lessThan: false}
var GreaterThanOrEqual = equalityCmp{lessThan: false, allowEqual: true}

type equalityCmp struct {
	allowEqual bool
	lessThan   bool
}

func (equalityCmp) ArgNames() []string {
	return []string{"got", "want"}
}

func (c equalityCmp) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	defer func() {
		// A panic is raised when the provided values are not comparable.
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	want := args[0]
	// Customize error message for non-nil errors.
	if _, ok := got.(error); ok && want == nil {
		return errors.New("got non-nil error")
	}

	gotType := reflect.TypeOf(got)
	wantType := reflect.TypeOf(want)
	if !c.comparable(gotType) || !c.comparable(wantType) {
		note("got type", qt.Unquoted(gotType.String()))
		note("want type", qt.Unquoted(wantType.String()))
		return errors.New("values are not comparable")
	}
	gotFloat := got.(float64)
	wantFloat := want.(float64)

	var result bool
	var cmpStr string
	switch {
	case c.lessThan && c.allowEqual:
		result = gotFloat <= wantFloat
		cmpStr = "<="
	case c.lessThan && !c.allowEqual:
		result = gotFloat < wantFloat
		cmpStr = "<"
	case !c.lessThan && c.allowEqual:
		result = gotFloat >= wantFloat
		cmpStr = ">="
	case !c.lessThan && !c.allowEqual:
		result = gotFloat > wantFloat
		cmpStr = ">"
	}
	if result {
		return nil
	}

	return fmt.Errorf("value %q is not %s %q", got, cmpStr, want)
}

func (c equalityCmp) comparable(typ reflect.Type) bool {
	var result float64
	assignType := reflect.ValueOf(result).Type()
	return typ.AssignableTo(assignType)
}
