package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	type Value struct {
		Name string
	}

	c := New[string, Value](time.Second * 2)
	c.Set("user1", Value{Name: "Alice"}, time.Second*1)
	c.Set("user2", Value{Name: "Bob"}, time.Second*5)
	c.Set("user3", Value{Name: "Eve"}, time.Second*10)
	c.Set("user4", Value{Name: "Charlie"}, NoExpiration)

	// Check that "user1" exists
	v, ok := c.Get("user1")
	if !ok || v.Name != "Alice" {
		t.Errorf("(Error) Expected to get 'Alice', but got %v", v.Name)
	}

	// Check that "user2" exists
	v, ok = c.Get("user2")
	if !ok || v.Name != "Bob" {
		t.Errorf("(Error) Expected to get 'Bob', but got %v", v.Name)
	}

	// Wait for 2 seconds
	time.Sleep(time.Second * 2)

	// Check that "user1" is expired
	_, ok = c.Get("user1")
	if ok {
		t.Error("(Error) Expected 'user1' to be expired and not found")
	}

	// Check that "user2" still exists
	v, ok = c.Get("user2")
	if !ok || v.Name != "Bob" {
		t.Errorf("(Error) Expected to get 'Bob', but got %v", v.Name)
	}

	// Test the Pop method
	poppedVal, popped := c.Pop("user2")
	if !popped || poppedVal.Name != "Bob" {
		t.Errorf("(Error) Expected to pop 'Bob', but got %v", poppedVal.Name)
	}

	// Check that "user2" no longer exists after Pop
	_, ok = c.Get("user2")
	if ok {
		t.Error("(Error) Expected 'user2' to be deleted after Pop")
	}

	// Test deleting an item
	c.Delete("user3")
	_, ok = c.Get("user3")
	if ok {
		t.Error("(Error) Expected 'user3' to be deleted")
	}

	v, ok = c.Get("user4")
	if !ok || v.Name != "Charlie" {
		t.Errorf("Expected to get 'Charlie', but got %v", v.Name)
	}

	return
}
