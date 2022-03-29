package cache

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	hash := NewMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	hash.Add("6", "4", "2") // 06 16 26 04 14 24 02 12 22； sorted：02 04 06 12 14 16 22 24 26

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"13": "4",
		"15": "6",
		"25": "6",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8") // 06 16 26 04 14 24 02 12 22 08 18 28； sorted：02 04 06 08 12 14 16 18 22 24 26 28

	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
