package shrink

import "testing"
import "testing/quick"

func TestShrink(t *testing.T) {
	t.Parallel()

	t.Run("shrink with one argument", func(t *testing.T) {
		fn := func(a int) bool {
			b := a % 1000
			return b < 500
		}
		if err := quick.Check(fn, nil); err != nil {
			if e := Shrink(fn, err); e != nil {
				err = e
			}
			t.Error(err)
		}
	})

	t.Run("shrink with two arguments", func(t *testing.T) {
		fn := func(a, c int) bool {
			b := a % 1000
			return b < 500
		}
		if err := quick.Check(fn, nil); err != nil {
			if e := Shrink(fn, err); e != nil {
				err = e
			}
			t.Error(err)
		}
	})
}
