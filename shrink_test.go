package shrink

import (
	"reflect"
	"testing"
	"testing/quick"
)

type IntAliasType int

func (x IntAliasType) Shrink() (reflect.Value, error) {
	return reflect.ValueOf(IntAliasType(x / 2)), nil
}

func TestShrink(t *testing.T) {
	t.Parallel()

	t.Run("shrink with one argument", func(t *testing.T) {
		fn := func(a int) bool {
			return a < 10
		}
		err := Shrink(fn, nil, &quick.CheckError{
			Count: 1,
			In: []interface{}{
				1000,
			},
		})

		chkErr, ok := err.(*CheckError)
		if !ok {
			t.Fatalf("unexpected CheckError, %v", err)
		}

		if expected, actual := 8, chkErr.Count; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := 1000, chkErr.In[0]; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := 7, chkErr.Succeeded[0]; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("shrink with multiple arguments", func(t *testing.T) {
		fn := func(a int, b string) bool {
			return a < 10
		}
		err := Shrink(fn, nil, &quick.CheckError{
			Count: 1,
			In: []interface{}{
				1000,
				"asd",
			},
		})

		chkErr, ok := err.(*CheckError)
		if !ok {
			t.Fatalf("unexpected CheckError, %v", err)
		}

		if expected, actual := 8, chkErr.Count; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := []interface{}{1000, "asd"}, chkErr.In; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []interface{}{7, ""}, chkErr.Succeeded; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("shrink with shrinkable arguments", func(t *testing.T) {
		fn := func(a IntAliasType, b string) bool {
			return a < 10
		}
		err := Shrink(fn, nil, &quick.CheckError{
			Count: 1,
			In: []interface{}{
				IntAliasType(1000),
				"asd",
			},
		})

		chkErr, ok := err.(*CheckError)
		if !ok {
			t.Fatalf("unexpected CheckError, %v", err)
		}

		if expected, actual := 8, chkErr.Count; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := []interface{}{IntAliasType(1000), "asd"}, chkErr.In; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []interface{}{IntAliasType(7), ""}, chkErr.Succeeded; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
