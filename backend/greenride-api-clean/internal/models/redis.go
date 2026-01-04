package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"greenride/internal/config"

	"github.com/go-redis/redis/v8"
)

var (
	Redis *redis.Client
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	// 使用DSN格式连接Redis
	dsn := cfg.Redis.Dsn
	if dsn == "" {
		return fmt.Errorf("redis DSN is empty")
	}

	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse Redis DSN: %w", err)
	}

	Redis = redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v. Continuing without Redis.", err)
		// 不返回错误，允许系统在没有Redis的情况下继续运行
	}

	return nil
}

// GetRedis 获取Redis客户端实例
func GetRedis() *redis.Client {
	return Redis
}

// 基础缓存操作
func SetCache(key string, value interface{}, expiration time.Duration) error {
	if Redis == nil {
		return nil // 开发环境可能没有Redis
	}
	return Redis.Set(context.Background(), key, value, expiration).Err()
}

func GetCache(key string) (string, error) {
	if Redis == nil {
		return "", errors.New("redis client not initialized")
	}
	return Redis.Get(context.Background(), key).Result()
}

// 简化别名
func Set(key string, value string, expiration time.Duration) error {
	return SetCache(key, value, expiration)
}

func Get(key string) (string, error) {
	return GetCache(key)
}

func SetInt(key string, value int64, expiration time.Duration) error {
	if Redis == nil {
		return nil
	}
	return Redis.Set(context.Background(), key, value, expiration).Err()
}

func GetInt(key string) (int64, error) {
	if Redis == nil {
		return 0, errors.New("redis client not initialized")
	}
	return Redis.Get(context.Background(), key).Int64()
}

func Delete(key string) error {
	if Redis == nil {
		return nil
	}
	return Redis.Del(context.Background(), key).Err()
}

// DelCache 删除缓存，与Delete功能相同，保留兼容性
func DelCache(keys ...string) error {
	if Redis == nil {
		return nil
	}
	return Redis.Del(context.Background(), keys...).Err()
}

// Exists 检查键是否存在
func Exists(key string) bool {
	if Redis == nil {
		return false
	}
	result, err := Redis.Exists(context.Background(), key).Result()
	return err == nil && result > 0
}

// ExistsCache 检查键是否存在，返回更详细的结果
func ExistsCache(keys ...string) (int64, error) {
	if Redis == nil {
		return 0, nil
	}
	return Redis.Exists(context.Background(), keys...).Result()
}

// 对象缓存操作
func SetObjectCache(key string, obj interface{}, expiration time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return SetCache(key, string(data), expiration)
}

func GetObjectCache(key string, obj interface{}) error {
	data, err := GetCache(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), obj)
}

// SetObjectToCache 与SetObjectCache功能相同，保留兼容性
func SetObjectToCache(key string, obj interface{}, expiration time.Duration) error {
	return SetObjectCache(key, obj, expiration)
}

// GetObjectFromCache 从缓存获取对象并反序列化
func GetObjectFromCache[T any](key string) (*T, error) {
	data, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	var obj T
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

// BatchGetObjectCache 批量获取对象缓存
func BatchGetObjectCache(keys []string, objType any) []any {
	if len(keys) == 0 {
		return make([]any, 0)
	}

	pipeline := Redis.Pipeline()
	commands := make(map[string]*redis.StringCmd)

	// 批量发送GET命令
	ctx := context.Background()
	for _, key := range keys {
		commands[key] = pipeline.Get(ctx, key)
	}

	// 执行批量操作
	_, err := pipeline.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Printf("BatchGetObjectCache pipeline.Exec error: %v", err)
		return nil
	}

	// 解析结果
	var result []any
	for _, cmd := range commands {
		jsonData, err := cmd.Result()
		if err != nil || jsonData == "" {
			continue
		}

		// 创建目标类型的新实例
		obj := reflect.New(reflect.TypeOf(objType).Elem()).Interface()
		err = json.Unmarshal([]byte(jsonData), obj)
		if err != nil {
			continue
		}

		result = append(result, obj)
	}

	return result
}

// CacheBuilder 缓存构建器函数类型
type CacheBuilder[T any] func() (*T, error)

// GetOrSetCache 通用的获取或设置缓存函数
func GetOrSetCache[T any](key string, ttl int, builder CacheBuilder[T]) (*T, error) {
	// 尝试从缓存获取
	if cached, err := GetObjectFromCache[T](key); err == nil {
		return cached, nil
	}

	// 缓存未命中，使用构建器获取数据
	obj, err := builder()
	if err != nil {
		return nil, err
	}

	// 将数据存入缓存
	if err := SetObjectToCache(key, obj, time.Duration(ttl)*time.Second); err != nil {
		fmt.Printf("Failed to set cache for key %s: %v\n", key, err)
	}

	return obj, nil
}

// DeleteCachePattern 删除匹配模式的所有缓存
func DeleteCachePattern(pattern string) error {
	if Redis == nil {
		return nil
	}

	ctx := context.Background()
	iter := Redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := Redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// 司机位置相关方法
func SetDriverLocation(driverID string, lat, lng float64, expiration time.Duration) error {
	key := fmt.Sprintf("driver_location:%s", driverID)
	value := fmt.Sprintf("%f,%f", lat, lng)
	return SetCache(key, value, expiration)
}

func GetDriverLocation(driverID string) (lat, lng float64, err error) {
	key := fmt.Sprintf("driver_location:%s", driverID)
	value, err := GetCache(key)
	if err != nil {
		return 0, 0, err
	}
	_, err = fmt.Sscanf(value, "%f,%f", &lat, &lng)
	return lat, lng, err
}

func DelDriverLocation(driverID string) error {
	key := fmt.Sprintf("driver_location:%s", driverID)
	return DelCache(key)
}

// 会话相关方法
func SetSession(sessionID string, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return SetCache(key, userID, expiration)
}

func GetSession(sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return GetCache(key)
}

func DelSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return DelCache(key)
}

// 限流相关方法
func IncrWithExpire(key string, expiration time.Duration) (int64, error) {
	ctx := context.Background()
	pipe := Redis.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// 分布式锁相关方法
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return Redis.SetNX(context.Background(), key, value, expiration).Result()
}

// FormatCacheKey 格式化缓存键
func FormatCacheKey(template string, args ...interface{}) string {
	return fmt.Sprintf("greenride:%s", fmt.Sprintf(template, args...))
}
