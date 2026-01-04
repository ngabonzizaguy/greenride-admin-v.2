package task

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var taskLock *redsync.Redsync

// InitLock 初始化任务锁
func InitLock(redisClient *redis.Client) {
	pool := goredis.NewPool(redisClient)
	taskLock = redsync.New(pool)
}

// Lock 获取任务锁
func Lock(taskID string) (*redsync.Mutex, error) {
	mutex := taskLock.NewMutex(
		"task_lock:"+taskID,
		redsync.WithExpiry(time.Minute), // 锁过期时间1分钟
		redsync.WithTries(1),            // 仅尝试一次
	)

	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	return mutex, nil
}

// Unlock 释放任务锁
func Unlock(mutex *redsync.Mutex) error {
	_, err := mutex.Unlock()
	return err
}
