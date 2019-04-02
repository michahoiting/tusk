package marshal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	yaml "gopkg.in/yaml.v2"
)

func TestMapList(t *testing.T) {
	s1 := []byte(`[a, b, c]`)
	s2 := []byte(`{a: a, b: b, c: c}`)

	lm1 := ListMap{}
	lm2 := ListMap{}

	if err := yaml.Unmarshal(s1, &lm1); err != nil {
		t.Fatal(err)
	}

	if err := yaml.Unmarshal(s2, &lm2); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(lm1, lm2) {
		t.Errorf("listmaps do not match:\n%s", cmp.Diff(lm1, lm2))
	}
}

func TestMapList_values(t *testing.T) {
	s := []byte(`{a: first, b: second, c: third}`)
	lm := ListMap{}

	if err := yaml.Unmarshal(s, &lm); err != nil {
		t.Fatal(err)
	}

	expected := ListMap{
		"a": "first",
		"b": "second",
		"c": "third",
	}

	if !cmp.Equal(lm, expected) {
		t.Errorf("listmaps do not match:\n%s", cmp.Diff(lm, expected))
	}
}
