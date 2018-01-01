package shrink

import (
	"fmt"
	"math"
	"reflect"
	"testing/quick"
)

// MaxRetries is the amount of times we want to attempt to shrink the arguments
var MaxRetries = 100

// ErrNoShrinkValue declares if the shrinking process can not generate an
// argument for the callable
var ErrNoShrinkValue = fmt.Errorf("invalid shrink value")

type shrinkable interface {
	// Dismantle a value from a concrete value to a reflection value
	// If an error is thrown then we exit immediately
	Dismantle() (reflect.Value, error)
}

// Shrink values for utilizing against checking values.
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

		case bool:
			x = !r

		case int:
			x = r / 2
		case int8:
			x = r / 2
		case int16:
			x = r / 2
		case int32:
			x = r / 2
		case int64:
			x = r / 2

		case uint:
			x = r / 2
		case uint8:
			x = r / 2
		case uint16:
			x = r / 2
		case uint32:
			x = r / 2
		case uint64:
			x = r / 2

		case float32:
			z := r / 2
			if z < 0 {
				x = float32(math.Ceil(float64(z)))
			} else {
				x = float32(math.Floor(float64(z)))
			}
		case float64:
			z := r / 2
			if z < 0 {
				x = math.Ceil(z)
			} else {
				x = math.Floor(z)
			}
		}

		res[k] = x
	}
	return res, nil
}
