package colorful

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

func TestHexColor(t *testing.T) {
	for _, tc := range []struct {
		hc HexColor
		s  string
	}{
		{HexColor{R: 0, G: 0, B: 0}, "#000000"},
		{HexColor{R: 1, G: 0, B: 0}, "#ff0000"},
		{HexColor{R: 0, G: 1, B: 0}, "#00ff00"},
		{HexColor{R: 0, G: 0, B: 1}, "#0000ff"},
		{HexColor{R: 1, G: 1, B: 1}, "#ffffff"},
	} {
		var gotHC HexColor
		if err := gotHC.Scan(tc.s); err != nil {
			t.Errorf("_.Scan(%q) == %v, want <nil>", tc.s, err)
		}
		if !reflect.DeepEqual(gotHC, tc.hc) {
			t.Errorf("_.Scan(%q) wrote %v, want %v", tc.s, gotHC, tc.hc)
		}
		if gotValue, err := tc.hc.Value(); err != nil || !reflect.DeepEqual(gotValue, tc.s) {
			t.Errorf("%v.Value() == %v, %v, want %v, <nil>", tc.hc, gotValue, err, tc.s)
		}
		gotString := tc.hc.String()
		if !reflect.DeepEqual(gotString, tc.s) {
			t.Errorf("_.String() == %v, want %v", gotString, tc.s)
		}
	}
}

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

func TestHexColorDriverValue(t *testing.T) {
	// database/sql passes query arguments by value, so HexColor (not just
	// *HexColor) has to implement driver.Valuer.
	got, err := driver.DefaultParameterConverter.ConvertValue(HexColor{R: 1, G: 0, B: 0})
	if err != nil || got != "#ff0000" {
		t.Errorf("ConvertValue(HexColor) == %v, %v, want #ff0000, <nil>", got, err)
	}
}
