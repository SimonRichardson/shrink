package shrink

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing/quick"
)

const maxRetries = 100

// ErrNoShrinkValue declares if the shrinking process can not generate an
// argument for the callable
var ErrNoShrinkValue = fmt.Errorf("invalid shrink value")

// A Shrinkable can shrink/reduce of its own type.
type Shrinkable interface {
	// Shrink a value from a concrete value to a reflection value
	// If an error is thrown then we exit immediately
	Shrink(interface{}) (reflect.Value, error)
}

// A Config structure contains options for running a test.
type Config struct {
	// MaxRetries sets the maximum number of iterations.
	MaxRetries int
	// CheckConfig is used for the configuration of quick check
	CheckConfig *quick.Config
}

var defaultConfig Config

// Check looks for an input to f, any function that returns bool,
// such that f returns false. It calls f repeatedly, with arbitrary
// values for each argument. If f returns false on a given input,
// Check attempts to shrink the input till the test case succeeds, failing that
// it will then return that input as a *CheckError.
func Check(fn interface{}, config *Config) error {
	if config == nil {
		config = &defaultConfig
		config.MaxRetries = maxRetries
	}

	err := quick.Check(fn, config.CheckConfig)
	if err != nil {
		if e := Shrink(fn, config, err); e != nil {
			err = e
		}
	}

	return err
}

// A CheckError is the result of Check finding an error.
type CheckError struct {
	Count     int
	In        []interface{}
	Succeeded []interface{}
}

func (s *CheckError) Error() string {
	return fmt.Sprintf("#%d: failed on input %s, but succeeded with %s", s.Count, toString(s.In), toString(s.Succeeded))
}

// Shrink values for utilizing against checking values.
func Shrink(fn interface{}, config *Config, err error) error {
	if chkErr, ok := err.(*quick.CheckError); ok {
		if config == nil {
			config = &defaultConfig
			config.MaxRetries = maxRetries
		}

		// Shrink the in input that failed
		var (
			origin          = chkErr.In
			fVal, fType, ok = functionAndType(fn)
		)
		if !ok {
			return quick.SetupError("argument is not a function")
		}

		if fType.NumOut() != 1 {
			return quick.SetupError("function does not return one value")
		}
		if fType.Out(0).Kind() != reflect.Bool {
			return quick.SetupError("function does not return a bool")
		}

		for i := 0; i < config.MaxRetries; i++ {
			values, snkErr := shrink(chkErr.In)
			if snkErr != nil {
				return snkErr
			}
			if len(values) != fType.NumIn() {
				return quick.SetupError("functions have different types")
			}
			if fVal.Call(values)[0].Bool() {
				// Report out what worked
				return &CheckError{chkErr.Count + 1, origin, toInterfaces(values)}
			}

			chkErr = &quick.CheckError{
				Count: chkErr.Count + 1,
				In:    toInterfaces(values),
			}
		}
	}
	return err
}

func functionAndType(f interface{}) (v reflect.Value, t reflect.Type, ok bool) {
	v = reflect.ValueOf(f)
	ok = v.Kind() == reflect.Func
	if !ok {
		return
	}
	t = v.Type()
	return
}

func toInterfaces(values []reflect.Value) []interface{} {
	ret := make([]interface{}, len(values))
	for i, v := range values {
		ret[i] = v.Interface()
	}
	return ret
}

func toString(interfaces []interface{}) string {
	s := make([]string, len(interfaces))
	for i, v := range interfaces {
		s[i] = fmt.Sprintf("%#v", v)
	}
	return strings.Join(s, ", ")
}

func shrink(args []interface{}) ([]reflect.Value, error) {
	res := make([]reflect.Value, len(args))
	for k, v := range args {
		var x interface{}

		switch r := v.(type) {
		case Shrinkable:
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
			if n := r / 2; n < 0 {
				x = float32(math.Ceil(float64(n)))
			} else {
				x = float32(math.Floor(float64(n)))
			}
		case float64:
			if n := r / 2; n < 0 {
				x = math.Ceil(n)
			} else {
				x = math.Floor(n)
			}

		case string:
			x = r[:(len(r)+1)/2]

		default:
			return nil, quick.SetupError(fmt.Sprintf("cannot create shrink value of type %T for argument %d", v, k))
		}

		res[k] = reflect.ValueOf(x)
	}
	return res, nil
}
