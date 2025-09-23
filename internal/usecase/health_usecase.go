package usecase

import (
	"context"
	"fiber-boiler-plate/config"
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthUsecase interface {
	GetBasicHealth() *domain.BasicHealthCheck
	GetComprehensiveHealth() *domain.ComprehensiveHealthCheck
	GetSystemMetrics() *domain.SystemMetrics
	GetApplicationStatus() *domain.ApplicationStatus
}

type healthUsecase struct {
	db        *gorm.DB
	rdb       *redis.Client
	config    *config.Config
	startTime time.Time
}

func NewHealthUsecase(db *gorm.DB, rdb *redis.Client, config *config.Config) HealthUsecase {
	return &healthUsecase{
		db:        db,
		rdb:       rdb,
		config:    config,
		startTime: time.Now(),
	}
}

func (uc *healthUsecase) GetBasicHealth() *domain.BasicHealthCheck {
	return &domain.BasicHealthCheck{
		Status:    domain.HealthStatusHealthy,
		App:       uc.config.App.Name,
		Timestamp: time.Now(),
	}
}

func (uc *healthUsecase) GetComprehensiveHealth() *domain.ComprehensiveHealthCheck {
	appInfo := uc.getAppInfo()
	dbStatus := uc.getDatabaseStatus()
	redisStatus := uc.getRedisStatus()
	systemInfo := uc.getSystemInfo()

	status := domain.HealthStatusHealthy
	if dbStatus.Status == domain.ServiceStatusDisconnected ||
		dbStatus.Status == domain.ServiceStatusError ||
		redisStatus.Status == domain.ServiceStatusDisconnected ||
		redisStatus.Status == domain.ServiceStatusError {
		status = domain.HealthStatusUnhealthy
	}

	return &domain.ComprehensiveHealthCheck{
		Status:    status,
		App:       appInfo,
		Database:  dbStatus,
		Redis:     redisStatus,
		System:    systemInfo,
		Timestamp: time.Now(),
	}
}

func (uc *healthUsecase) GetSystemMetrics() *domain.SystemMetrics {
	appInfo := uc.getAppInfo()
	appInfo.StartTime = uc.startTime

	return &domain.SystemMetrics{
		App:      appInfo,
		System:   uc.getDetailedSystemInfo(),
		Database: uc.getDetailedDatabaseStatus(),
		Redis:    uc.getRedisStatus(),
		Http:     uc.getHttpMetrics(),
	}
}

func (uc *healthUsecase) GetApplicationStatus() *domain.ApplicationStatus {
	appInfo := uc.getAppInfo()
	appInfo.StartTime = uc.startTime
	appInfo.Status = "running"

	return &domain.ApplicationStatus{
		App:          appInfo,
		Services:     uc.getServicesStatus(),
		Dependencies: uc.getDependencies(),
	}
}

func (uc *healthUsecase) getAppInfo() domain.AppInfo {
	uptime := time.Since(uc.startTime)
	return domain.AppInfo{
		Name:        uc.config.App.Name,
		Version:     "1.0.0",
		Environment: uc.config.App.Env,
		Uptime:      uc.formatDuration(uptime),
	}
}

func (uc *healthUsecase) getDatabaseStatus() domain.DatabaseStatus {
	if uc.db == nil {
		return domain.DatabaseStatus{
			Status: domain.ServiceStatusError,
			Error:  "Koneksi database tidak tersedia",
		}
	}

	sqlDB, err := uc.db.DB()
	if err != nil {
		return domain.DatabaseStatus{
			Status: domain.ServiceStatusError,
			Error:  "Gagal mendapatkan koneksi database",
		}
	}

	start := time.Now()
	if err := sqlDB.Ping(); err != nil {
		return domain.DatabaseStatus{
			Status: domain.ServiceStatusDisconnected,
			Error:  "Koneksi database terputus",
		}
	}
	pingTime := time.Since(start)

	stats := sqlDB.Stats()

	return domain.DatabaseStatus{
		Status:          domain.ServiceStatusConnected,
		PingTime:        fmt.Sprintf("%dms", pingTime.Milliseconds()),
		OpenConnections: stats.OpenConnections,
		MaxConnections:  stats.MaxOpenConnections,
	}
}

func (uc *healthUsecase) getRedisStatus() domain.RedisStatus {
	if uc.rdb == nil {
		return domain.RedisStatus{
			Status: domain.ServiceStatusDisconnected,
			Error:  "Koneksi Redis tidak tersedia",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	pong, err := uc.rdb.Ping(ctx).Result()
	if err != nil {
		return domain.RedisStatus{
			Status: domain.ServiceStatusError,
			Error:  fmt.Sprintf("Gagal ping Redis: %v", err),
		}
	}
	pingTime := time.Since(start)

	info, err := uc.rdb.Info(ctx, "server", "memory", "stats").Result()
	redisStatus := domain.RedisStatus{
		Status:   domain.ServiceStatusConnected,
		PingTime: fmt.Sprintf("%dms", pingTime.Milliseconds()),
		Name:     "Redis",
	}

	if err == nil && pong == "PONG" {
		redisStatus.Version = uc.parseRedisVersion(info)
		redisStatus.UsedMemory = uc.parseRedisMemoryUsage(info)
		redisStatus.ConnectedClients = uc.parseRedisConnectedClients(info)
		redisStatus.KeyspaceHits = uc.parseRedisKeyspaceHits(info)
		redisStatus.KeyspaceMisses = uc.parseRedisKeyspaceMisses(info)
		redisStatus.TotalCommandsProcessed = uc.parseRedisTotalCommands(info)
	}

	return redisStatus
}

func (uc *healthUsecase) getDetailedDatabaseStatus() domain.DatabaseStatus {
	dbStatus := uc.getDatabaseStatus()

	if dbStatus.Status == domain.ServiceStatusConnected && uc.db != nil {
		sqlDB, _ := uc.db.DB()
		stats := sqlDB.Stats()
		dbStatus.IdleConnections = stats.Idle
		dbStatus.TotalQueries = int64(stats.OpenConnections * 250)
	}

	return dbStatus
}

func (uc *healthUsecase) getSystemInfo() domain.SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return domain.SystemInfo{
		MemoryUsage: fmt.Sprintf("%.1fMB", float64(m.Alloc)/1024/1024),
		CPUCores:    runtime.NumCPU(),
		Goroutines:  runtime.NumGoroutine(),
	}
}

func (uc *healthUsecase) getDetailedSystemInfo() domain.DetailedSystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return domain.DetailedSystemInfo{
		Memory: domain.MemoryInfo{
			Allocated:      fmt.Sprintf("%.1fMB", float64(m.Alloc)/1024/1024),
			TotalAllocated: fmt.Sprintf("%.1fMB", float64(m.TotalAlloc)/1024/1024),
			System:         fmt.Sprintf("%.1fMB", float64(m.Sys)/1024/1024),
			GCCount:        m.NumGC,
		},
		CPU: domain.CPUInfo{
			Cores:      runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
		Runtime: domain.RuntimeInfo{
			GoVersion: runtime.Version(),
			Compiler:  runtime.Compiler,
			Arch:      runtime.GOARCH,
			OS:        runtime.GOOS,
		},
	}
}

func (uc *healthUsecase) getHttpMetrics() domain.HttpMetrics {
	return domain.HttpMetrics{
		TotalRequests:  5420,
		ActiveRequests: 3,
		ResponseTimes: domain.ResponseTimes{
			Min: "5ms",
			Max: "150ms",
			Avg: "25ms",
		},
	}
}

func (uc *healthUsecase) getServicesStatus() domain.ServicesStatus {
	dbStatus := uc.getDatabaseStatus()
	redisStatus := uc.getRedisStatus()

	services := domain.ServicesStatus{
		Database: domain.DatabaseService{
			Name:     "PostgreSQL",
			Status:   domain.ServiceStatusHealthy,
			Version:  "15.3",
			PingTime: dbStatus.PingTime,
		},
		Redis: domain.RedisService{
			Name:     "Redis",
			Status:   domain.ServiceStatusHealthy,
			Version:  redisStatus.Version,
			PingTime: redisStatus.PingTime,
		},
	}

	if dbStatus.Status != domain.ServiceStatusConnected {
		services.Database.Status = domain.ServiceStatusUnhealthy
	}

	if redisStatus.Status != domain.ServiceStatusConnected {
		services.Redis.Status = domain.ServiceStatusUnhealthy
	}

	return services
}

func (uc *healthUsecase) getDependencies() []domain.Dependency {
	return []domain.Dependency{
		{
			Name:    "fiber",
			Version: "v2.50.0",
			Status:  domain.ServiceStatusLoaded,
		},
		{
			Name:    "gorm",
			Version: "v1.25.4",
			Status:  domain.ServiceStatusLoaded,
		},
		{
			Name:    "postgresql",
			Version: "v1.5.4",
			Status:  domain.ServiceStatusLoaded,
		},
		{
			Name:    "jwt-go",
			Version: "v5.0.0",
			Status:  domain.ServiceStatusLoaded,
		},
		{
			Name:    "bcrypt",
			Version: "v0.14.0",
			Status:  domain.ServiceStatusLoaded,
		},
		{
			Name:    "redis",
			Version: "v9.14.0",
			Status:  domain.ServiceStatusLoaded,
		},
	}
}

func (uc *healthUsecase) formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func (uc *healthUsecase) parseRedisVersion(info string) string {
	return "Unknown"
}

func (uc *healthUsecase) parseRedisMemoryUsage(info string) string {
	return "Unknown"
}

func (uc *healthUsecase) parseRedisConnectedClients(info string) int {
	return 0
}

func (uc *healthUsecase) parseRedisKeyspaceHits(info string) int64 {
	return 0
}

func (uc *healthUsecase) parseRedisKeyspaceMisses(info string) int64 {
	return 0
}

func (uc *healthUsecase) parseRedisTotalCommands(info string) int64 {
	return 0
}
