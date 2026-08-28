package colorful

import (
	"encoding/json"
	"reflect"
	"testing"
)

type CompositeType struct {
	Name  string   `json:"name,omitempty"`
	Color HexColor `json:"color,omitempty"`
}

func TestHexColorCompositeJson(t *testing.T) {
	var obj = CompositeType{Name: "John", Color: HexColor{R: 1, G: 0, B: 1}}
	var jsonData, err = json.Marshal(obj)
	if err != nil {
		t.Errorf("json.Marshall(obj) wrote %v", err)
	}
	var obj2 CompositeType
	err = json.Unmarshal(jsonData, &obj2)

	if err != nil {
		t.Errorf("json.Unmarshall(%s) wrote %v", jsonData, err)
	}

	if !reflect.DeepEqual(obj2, obj) {
		t.Errorf("json.Unmarshal(json.Marsrhall(obj)) != obj")
	}

}

// TestHexColorString covers String on every target, including the ones where
// the database/sql methods are built out.
func TestHexColorString(t *testing.T) {
	for _, tc := range []struct {
		hc HexColor
		s  string
	}{
		{HexColor{R: 0, G: 0, B: 0}, "#000000"},
		{HexColor{R: 1, G: 0, B: 1}, "#ff00ff"},
		{HexColor{R: 1, G: 1, B: 1}, "#ffffff"},
	} {
		if got := tc.hc.String(); got != tc.s {
			t.Errorf("_.String() == %v, want %v", got, tc.s)
		}
	}
}
