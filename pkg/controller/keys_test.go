package controller

import (
  "testing"
)

func TestSplitOnce(t *testing.T) {
  // t.Log("Oh noes - something is false")
  // t.Fail()
  a,b := splitOnce("key:value", ":")
  if a != "key" {
    t.Fatalf("a should equal 'key' (got %v)", a)
  }
  if b != "value" {
    t.Fatalf("a should equal 'value' (got %v)", a)
  }

  return
}
