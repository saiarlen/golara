package constants

const (
	// Framework info
	FrameworkName    = "Golara"
	FrameworkVersion = "1.0.0"
	
	// Default configurations
	DefaultPort        = "9000"
	DefaultBodyLimit   = 20 * 1024 * 1024 // 20MB
	DefaultCachePrefix = "golara"
	DefaultQueueName   = "default"
	
	// Environment values
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
	
	// Cache drivers
	CacheDriverMemory = "memory"
	CacheDriverRedis  = "redis"
	
	// Queue drivers
	QueueDriverSync   = "sync"
	QueueDriverRedis  = "redis"
	QueueDriverMemory = "memory"
)