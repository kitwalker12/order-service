package interfaces

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	Uri string
}

func (cache Cache) Do(cmd string, args ...interface{}) (interface{}, error) {
	c, err := redis.DialURL(cache.Uri)
	if err != nil {
		return nil, fmt.Errorf("DialURL: %s", err)
	}
	defer c.Close()

	n, err := c.Do(cmd, args)
	if err != nil {
		return nil, fmt.Errorf("Do: %s", err)
	}

	return n, nil
}
