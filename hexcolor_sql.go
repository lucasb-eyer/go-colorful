//go:build !tinygo && !js

package colorful

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// HexColor's database/sql integration.
//
// It is built out on TinyGo and on js/wasm because importing database/sql/driver
// there is pure cost: driver.Value is the only thing this package needs from it,
// and that single return type pulls database/sql/driver -> uuid -> crypto/rand
// -> the whole crypto/fips140 tree. On a js/wasm build that is 26 crypto
// packages and ~722KB of a 3.2MB binary, none of which can run: neither target
// has a database/sql driver to be a Valuer for.
//
// Everywhere else this builds and behaves exactly as before.

type errUnsupportedType struct {
	got  interface{}
	want reflect.Type
}

func (hc *HexColor) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errUnsupportedType{got: reflect.TypeOf(value), want: reflect.TypeOf("")}
	}
	c, err := Hex(s)
	if err != nil {
		return err
	}
	*hc = HexColor(c)
	return nil
}

func (hc *HexColor) Value() (driver.Value, error) {
	return Color(*hc).Hex(), nil
}

func (e errUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type: got %v, want a %s", e.got, e.want)
}
