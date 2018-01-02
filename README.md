# shrink

Shrink is a drop in replacement for the stdlib `quick.Check`, but with the
additional benefits of shrinking the input arguments on failure.

## Example

```go
func TestDivide(t *testing.T) {
  fn := func(x int) bool {
    return (x % 1000) < 10
  }
  if err := shrink.Check(fn, nil); err != nil {
    t.Error(err)
  }
}
```
