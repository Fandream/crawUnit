package database

import (
	"context"
	"crawlab/entity"
	"crawlab/utils"
	"errors"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"runtime/debug"
	"strings"
	"time"
)

var RedisClient *Redis

type Redis struct {
	pool *redis.Pool
}

type Mutex struct {
	Name   string
	expiry time.Duration
	tries  int
	delay  time.Duration
	value  string
}

func NewRedisClient() *Redis {
	return &Redis{pool: NewRedisPool()}
}

func (r *Redis) RPush(collection string, value interface{}) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) LPush(collection string, value interface{}) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) LPop(collection string) (string, error) {
	c := r.pool.Get()
	defer utils.Close(c)

	value, err2 := redis.String(c.Do("LPOP", collection))
	if err2 != nil {
		return value, err2
	}
	return value, nil
}

func (r *Redis) HSet(collection string, key string, value string) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HSET", collection, key, value); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
func (r *Redis) Ping() error {
	c := r.pool.Get()
	defer utils.Close(c)
	_, err2 := redis.String(c.Do("PING"))
	return err2
}
func (r *Redis) HGet(collection string, key string) (string, error) {
	c := r.pool.Get()
	defer utils.Close(c)
	value, err2 := redis.String(c.Do("HGET", collection, key))
	if err2 != nil && err2 != redis.ErrNil {
		log.Error(err2.Error())
		debug.PrintStack()
		return value, err2
	}
	return value, nil
}

func (r *Redis) HDel(collection string, key string) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HDEL", collection, key); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
func (r *Redis) HScan(collection string) (results []string, err error) {
	c := r.pool.Get()
	defer utils.Close(c)
	var (
		cursor int64
		items  []string
	)

	for {
		values, err := redis.Values(c.Do("HSCAN", collection, cursor))
		if err != nil {
			return results, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			return results, err
		}

		results = append(results, items[1])

		if cursor == 0 {
			break
		}
	}
	return results, err

}
func (r *Redis) HKeys(collection string) ([]string, error) {
	c := r.pool.Get()
	defer utils.Close(c)

	value, err2 := redis.Strings(c.Do("HKEYS", collection))
	if err2 != nil {
		log.Error(err2.Error())
		debug.PrintStack()
		return []string{}, err2
	}
	return value, nil
}

func (r *Redis) BRPop(collection string, timeout int) (string, error) {
	if timeout <= 0 {
		timeout = 60
	}
	c := r.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BRPOP", collection, timeout))
	if err != nil {
		return "", err
	}
	return values[1], nil
}


