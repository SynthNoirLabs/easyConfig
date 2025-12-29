package cache

import (
	"testing"
	"time"
)

func TestCache_Set_Get(t *testing.T) {
	c := New()
	key := "test_key"
	value := "test_value"

	c.Set(key, value, 1*time.Second)

	retrievedValue, found, stale := c.Get(key)
	if !found {
		t.Error("Expected to find the key in the cache")
	}
	if stale {
		t.Error("Expected the key to not be stale")
	}
	if retrievedValue != value {
		t.Errorf("Expected value %s, but got %s", value, retrievedValue)
	}
}

func TestCache_TTL(t *testing.T) {
	c := New()
	key := "test_key"
	value := "test_value"

	c.Set(key, value, 1*time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	_, found, stale := c.Get(key)
	if !found {
		t.Error("Expected to find the key in the cache")
	}
	if !stale {
		t.Error("Expected the key to be stale")
	}
}

func TestCache_Delete(t *testing.T) {
	c := New()
	key := "test_key"
	value := "test_value"

	c.Set(key, value, 1*time.Second)
	c.Delete(key)

	_, found, _ := c.Get(key)
	if found {
		t.Error("Expected to not find the key in the cache after deletion")
	}
}
