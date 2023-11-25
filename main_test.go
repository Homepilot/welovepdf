package main

import (
	"log"
	"testing"
)

func TestFooer(t *testing.T) {
	result := Fooer(3)
	log.Println("Running test")
	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}
