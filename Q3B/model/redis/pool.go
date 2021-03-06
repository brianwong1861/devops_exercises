package  redis

import (
	"os/signal"
	"time"
	"github.com/garyburd/redigo/redis"
	"os"
	"syscall"
	"github.com/joho/godotenv"
	
	
)

var (
	Pool *redis.Pool
)

func init() {
	_ = godotenv.Load(".env")

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	
	Pool = newPool(redisHost)
	// Redis initializing test 
	err := Ping()
	if err != nil {
		panic("Redis initializing failed")
	} 
	cleanupHook()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error){
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func(){
		<- c
		Pool.Close()
		os.Exit(0)
	}()
}