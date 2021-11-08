package atomic

// String is an atomic type-safe wrapper for string values.
type String struct {
	_ nocmp // disallow non-atomic comparison

	v Value
}

var _zeroString string

// NewString creates a new String.
func NewString(val string) *String {
	x := &String{}
	if val != _zeroString {
		x.Store(val)
	}
	return x
}

// Load atomically loads the wrapped string.
func (x *String) Load() string {
	if v := x.v.Load(); v != nil {
		return v.(string)
	}
	return _zeroString
}

// Store atomically stores the passed string.
func (x *String) Store(val string) {
	x.v.Store(val)
}
