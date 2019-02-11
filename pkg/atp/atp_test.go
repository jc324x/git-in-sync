package atp

import (
	"reflect"
	"testing"
)

func TestResulter(t *testing.T) {
	want := Results{
		{"hendricius", "github", "recipes", []string{"pizza-dough"}},
		{"rochacbruno", "github", "recipes", []string{"vegan_recipes"}},
		{"niw", "github", "recipes", []string{"ramen"}},
	}

	got := Resulter("recipes")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Resulter: != DeepEqual (%v != %v)", got, want)
	}
}
