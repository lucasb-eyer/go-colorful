package colorful

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
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

func TestHexColorYaml(t *testing.T) {
	obj := HexColor{R: 0, G: 1, B: 0}
	want := "'#00ff00'\n"
	raw, err := yaml.Marshal(&obj)
	if err != nil {
		t.Error(err)
	}

	if strings.Compare(want, string(raw)) != 0 {
		t.Errorf("Wanted '%s' got '%s'.", want, string(raw))
	}

	var got HexColor
	err = yaml.Unmarshal(raw, &got)
	if err != nil {
		t.Error(err)
	}
	if obj != got {
		t.Error(got)
	}
}

func TestHexColorCompositeYaml(t *testing.T) {
	obj := CompositeType{Name: "John", Color: HexColor{R: 0, G: 1, B: 0}}
	raw, err := yaml.Marshal(&obj)
	if err != nil {
		t.Error(err)
	}

	var got CompositeType
	err = yaml.Unmarshal(raw, &got)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, obj) {
		t.Errorf("yaml.Unmarshal(yaml.Marshal(obj)) != obj")
	}
}
