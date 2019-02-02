package atp

import (
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	_, clean := Setup("atp", "recipes")
	defer clean()
}

func TestWant(t *testing.T) {
	want := Results{
		{"hendricius", "github", "recipes", []string{"pizza-dough", "the-bread-code"}},
		{"cocktails-for-programmers", "github", "recipes", []string{"cocktails-for-programmers"}},
		{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
		{"niw", "github", "recipes", []string{"ramen"}},
	}

	got := Resulter("recipes")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Resulter: != DeepEqual (%v != %v)", got, want)
	}
}
