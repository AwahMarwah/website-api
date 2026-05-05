package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func GenerateCacheKey(prefix string, query interface{}) string {
	data, err := json.Marshal(query)
	if err != nil {
		return fmt.Sprintf("%s:default", prefix)
	}

	h := sha1.New()
	h.Write(data)
	hash := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%s", prefix, hash)
}
