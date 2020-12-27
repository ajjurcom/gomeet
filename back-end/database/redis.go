package database

import "github.com/gomodule/redigo/redis"

var RedisDB redis.Conn

func InitRedis() (err error) {
	RedisDB, err = redis.Dial("tcp", "127.0.0.1:6379")
	return
}

func CloseRedis() {
	RedisDB.Close()
}
