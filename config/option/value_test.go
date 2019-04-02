package option

import (
	"reflect"
	"testing"

	"github.com/rliebz/tusk/config/marshal"
	yaml "gopkg.in/yaml.v2"
)

func TestValue_UnmarshalYAML(t *testing.T) {
	s1 := []byte(`value: example`)
	s2 := []byte(`example`)
	v1 := Value{}
	v2 := Value{}

	if err := yaml.Unmarshal(s1, &v1); err != nil {
		t.Fatalf("yaml.Unmarshal(%s, ...): unexpected error: %s", s1, err)
	}

	if err := yaml.Unmarshal(s2, &v2); err != nil {
		t.Fatalf("yaml.Unmarshal(%s, ...): unexpected error: %s", s2, err)
	}

	if !reflect.DeepEqual(v1, v2) {
		t.Errorf(
			"Unmarshalling of values `%s` and `%s` not equal:\n%#v != %#v",
			s1, s2, v1, v2,
		)
	}

	if v1.Value != "example" {
		t.Errorf(
			"yaml.Unmarshal(%s, ...): expected member `%s`, actual `%s`",
			s1, "example", v1.Command,
		)
	}
}

func TestValue_UnmarshalYAML_value_and_command(t *testing.T) {
	s := []byte(`{value: "example", command: "echo hello"}`)
	v := Value{}

	if err := yaml.Unmarshal(s, &v); err == nil {
		t.Fatalf(
			"yaml.Unmarshal(%s, ...): expected err, actual nil", s,
		)
	}
}

func TestValueList_UnmarshalYAML(t *testing.T) {
	s1 := []byte(`example`)
	s2 := []byte(`[example]`)
	v1 := ValueList{}
	v2 := ValueList{}

	if err := yaml.Unmarshal(s1, &v1); err != nil {
		t.Fatalf("yaml.Unmarshal(%s, ...): unexpcted error: %s", s1, err)
	}

	if err := yaml.Unmarshal(s2, &v2); err != nil {
		t.Fatalf("yaml.Unmarshal(%s, ...): unexpcted error: %s", s2, err)
	}

	if !reflect.DeepEqual(v1, v2) {
		t.Errorf(
			"Unmarshalling of valueLists `%s` and `%s` not equal:\n%#v != %#v",
			s1, s2, v1, v2,
		)
	}

	if len(v1) != 1 {
		t.Errorf(
			"yaml.Unmarshal(%s, ...): expected 1 item, actual %d",
			s1, len(v1),
		)
	}

	if v1[0].Value != "example" {
		t.Errorf(
			"yaml.Unmarshal(%s, ...): expected member `%s`, actual `%s`",
			s1, "example", v1[0].Value,
		)
	}
}

func TestValueWithList(t *testing.T) {
	cases := []struct {
		name    string
		values  marshal.ListMap
		value   string
		expect  string
		wantErr bool
	}{
		{
			"nil-values",
			nil,
			"foo",
			"foo",
			false,
		},
		{
			"no-values",
			marshal.ListMap{},
			"foo",
			"foo",
			false,
		},
		{
			"unmapped-value",
			marshal.ListMap{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},
			"foo",
			"foo",
			false,
		},
		{
			"mapped-value",
			marshal.ListMap{
				"foo": "foovalue",
				"bar": "barvalue",
				"baz": "bazvalue",
			},
			"foo",
			"foovalue",
			false,
		},
		{
			"bad-value",
			marshal.ListMap{
				"foo": "foovalue",
				"bar": "barvalue",
				"baz": "bazvalue",
			},
			"unknown",
			"",
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v := valueWithList{
				ValuesAllowed: tt.values,
			}

			actual, err := v.mappedValue(tt.value, "some-descriptor")
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("want error: %t, got %q", tt.wantErr, err)
			}

			if actual != tt.expect {
				t.Errorf("want value %q, got %q", tt.expect, actual)
			}
		})
	}
}
