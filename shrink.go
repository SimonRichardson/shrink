package shrink

import (
	"fmt"
	"reflect"
	"testing/quick"
)

var MaxRetries = 100

var ErrNoShrinkValue = fmt.Errorf("invalid shrink value")

type shrinkable interface {
	Dismantal() (reflect.Value, error)
}

func Shrink(fn interface{}, err error) error {
	if chkErr, ok := err.(*quick.CheckError); ok {
		// Shrink the in input that failed
		var (
			tries      = 0
			origin     = chkErr.In
			callable   = reflect.ValueOf(fn)
			callableIn = callable.Type().NumIn()
		)
		for {
			if tries == MaxRetries {
				return err
			}

			values, snkErr := shrink(chkErr.In)
			if snkErr != nil {
				return err
			}
			if len(values) != callableIn {
				return fmt.Errorf("shrink: function values missmatch")
			}

			result := callable.Call(refValue(values))

			if len(result) != 1 {
				return fmt.Errorf("expected bool result")
			}
			if result[0].Bool() {
				// Report out what worked
				return fmt.Errorf("shrink: failed with (%v), but succeeded with (%v)", origin, values)
			}

			tries++
		}
	}
	return err
}

func refValue(args []interface{}) []reflect.Value {
	res := make([]reflect.Value, len(args))
	for k, v := range args {
		res[k] = reflect.ValueOf(v)
	}
	return res
}

func shrink(args []interface{}) ([]interface{}, error) {
	res := make([]interface{}, len(args))
	for k, v := range args {
		var x interface{}

		switch r := v.(type) {
		case shrinkable:
			// TODO:

		case int:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case int8:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case int16:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case int32:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case int64:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}

		case uint:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case uint8:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case uint16:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case uint32:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		case uint64:
			if x = r / 2; x == 0 {
				return nil, ErrNoShrinkValue
			}
		}

		res[k] = x
	}
	return res, nil
}
