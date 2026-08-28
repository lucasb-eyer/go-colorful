package colorful

import (
	"encoding/json"
)

// A HexColor is a Color stored as a hex string "#rrggbb". It implements the
// encoding/json and gopkg.in/yaml marshaller interfaces, and the envconfig
// Decoder. The database/sql integration lives in hexcolor_sql.go, which is
// built only where a database/sql consumer can exist.
type HexColor Color

func (hc HexColor) String() string {
	return Color(hc).Hex()
}

func (hc *HexColor) UnmarshalJSON(data []byte) error {
	var hexCode string
	if err := json.Unmarshal(data, &hexCode); err != nil {
		return err
	}

	var col, err = Hex(hexCode)
	if err != nil {
		return err
	}
	*hc = HexColor(col)
	return nil
}

func (hc HexColor) MarshalJSON() ([]byte, error) {
	return json.Marshal(Color(hc).Hex())
}

// Decode - deserialize function for https://github.com/kelseyhightower/envconfig
func (hc *HexColor) Decode(hexCode string) error {
	var col, err = Hex(hexCode)
	if err != nil {
		return err
	}
	*hc = HexColor(col)
	return nil
}

func (hc HexColor) MarshalYAML() (interface{}, error) {
	return Color(hc).Hex(), nil
}

func (hc *HexColor) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var hexCode string
	if err := unmarshal(&hexCode); err != nil {
		return err
	}

	var col, err = Hex(hexCode)
	if err != nil {
		return err
	}

	*hc = HexColor(col)

	return nil
}
