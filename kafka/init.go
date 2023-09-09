package kafka 

import (
	"github.com/patrickmn/go-cache"
)

var (
	schemaCache cache
)

func init() {
	if cache == true {
		schemaCache := cache.New(2*time.Hour, 10*time.Minute)
	}
}