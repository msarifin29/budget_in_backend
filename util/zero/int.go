package zero

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int is a nullable int64.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int struct {
	sql.NullInt64
}

// Int64 is an alias for Int.
type Int64 = Int

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will be null if zero.
func IntFrom(i int64) Int {
	return NewInt(i, i != 0)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	n := NewInt(*i, true)
	return n
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int.
func (i *Int) UnmarshalJSON(data []byte) error {
	err := UnmarshalIntJSON(data, &i.Int64, &i.Valid, 64, strconv.ParseInt)
	if err != nil {
		return err
	}
	i.Valid = i.Int64 != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank, or zero.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	err := UnmarshalIntText(text, &i.Int64, &i.Valid, 64, strconv.ParseInt)
	if err != nil {
		return err
	}
	i.Valid = i.Int64 != 0
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for null or zero Ints, for future omitempty support (Go 1.4?)
func (i Int) IsZero() bool {
	return !i.Valid || i.Int64 == 0
}

// Equal returns true if both ints have the same value or are both either null or zero.
func (i Int) Equal(other Int) bool {
	return i.ValueOrZero() == other.ValueOrZero()
}

// func (i Int) value() (int64, bool) {
// 	return i.Int64, i.Valid
// }

type Integer interface {
	int64 | int32 | int16 | byte
}

func UnmarshalIntJSON[T Integer, U int64 | uint64](data []byte, value *T, valid *bool, bits int, parse func(string, int, int) (U, error)) error {
	if len(data) == 0 {
		return fmt.Errorf("UnmarshalJSON: no data")
	}

	switch data[0] {
	case 'n':
		*value = 0
		*valid = false
		return nil

	case '"':
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
		}
		n, err := parse(str, 10, bits)
		if err != nil {
			return fmt.Errorf("null: couldn't convert string to int: %w", err)
		}
		*value = T(n)
		*valid = true
		return nil

	default:
		err := json.Unmarshal(data, value)
		*valid = err == nil
		return err
	}
}

func UnmarshalIntText[T Integer, U int64 | uint64](text []byte, value *T, valid *bool, bits int, parse func(string, int, int) (U, error)) error {
	str := string(text)
	if str == "" || str == "null" {
		*value = 0
		*valid = false
		return nil
	}
	n, err := parse(str, 10, bits)
	*value = T(n)
	if err != nil {
		*valid = false
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	*valid = true
	return nil
}