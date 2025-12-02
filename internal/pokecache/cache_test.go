package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cases := map[string]struct {
		key string
		val []byte
	}{
		"get":     {"https://example.com", []byte("testdata")},
		"getSame": {"https://example.com", []byte("testdatacorrect")},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cache := NewCache(time.Second * 5)
			cache.Add(tc.key, tc.val)
			val, ok := cache.Get(tc.key)
			if !ok {
				t.Errorf("key not found")
				return
			}
			if string(val) != string(tc.val) {
				t.Errorf("value not found")
				return
			}
		})
	}

}

func TestReapLoop(t *testing.T) {
	const basetime = 200 * time.Millisecond
	const waittime = 200 * time.Millisecond

	cache := NewCache(basetime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waittime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
