package atp

import (
	"log"
	"testing"
)

func TestSetup(t *testing.T) {
	Setup("atp", "recipes")
	log.Print(Tmap)
}
