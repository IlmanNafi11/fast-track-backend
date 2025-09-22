package usecase

import (
	"fiber-boiler-plate/config"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(sql string, values ...interface{}) *gorm.DB {
	args := m.Called(sql, values)
	result := args.Get(0).(*gorm.DB)
	return result
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	result := args.Get(0).(*gorm.DB)
	return result
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	result := args.Get(0).(*gorm.DB)
	return result
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	m.Called(value, conds)
	return &gorm.DB{}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return &gorm.DB{}
}

func TestHealthUsecase_GetBasicHealth_Success(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
		},
	}

	healthUsecase := usecase.NewHealthUsecase(nil, nil, cfg)
	result := healthUsecase.GetBasicHealth()

	assert.NotNil(t, result)
	assert.Equal(t, domain.HealthStatusHealthy, result.Status)
	assert.Equal(t, "test-app", result.App)
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
}

func TestHealthUsecase_GetComprehensiveHealth_Success(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
		},
	}

	healthUsecase := usecase.NewHealthUsecase(nil, nil, cfg)
	result := healthUsecase.GetComprehensiveHealth()

	assert.NotNil(t, result)
	assert.Equal(t, "test-app", result.App.Name)
	assert.Equal(t, "1.0.0", result.App.Version)
	assert.Equal(t, "test", result.App.Environment)
	assert.NotEmpty(t, result.App.Uptime)
	assert.Equal(t, domain.HealthStatusUnhealthy, result.Status)

	// Test Redis status with nil client
	assert.NotNil(t, result.Redis)
	assert.Equal(t, domain.ServiceStatusDisconnected, result.Redis.Status)
	assert.Equal(t, "Koneksi Redis tidak tersedia", result.Redis.Error)
	assert.Equal(t, domain.ServiceStatusError, result.Database.Status)
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
}

func TestHealthUsecase_GetSystemMetrics_Success(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
		},
	}

	healthUsecase := usecase.NewHealthUsecase(nil, nil, cfg)
	result := healthUsecase.GetSystemMetrics()

	assert.NotNil(t, result)
	assert.Equal(t, "test-app", result.App.Name)
	assert.Equal(t, "1.0.0", result.App.Version)
	assert.Equal(t, "test", result.App.Environment)
	assert.NotEmpty(t, result.App.Uptime)
	assert.NotEmpty(t, result.System.Memory.Allocated)
	assert.Greater(t, result.System.CPU.Cores, 0)
	assert.GreaterOrEqual(t, result.System.CPU.Goroutines, 0)
	assert.NotEmpty(t, result.System.Runtime.GoVersion)
	assert.NotEmpty(t, result.System.Runtime.OS)
	assert.Greater(t, result.Http.TotalRequests, int64(0))
	assert.Equal(t, domain.ServiceStatusError, result.Database.Status)

	// Test Redis metrics with nil client
	assert.NotNil(t, result.Redis)
	assert.Equal(t, domain.ServiceStatusDisconnected, result.Redis.Status)
}

func TestHealthUsecase_GetApplicationStatus_Success(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
		},
	}

	healthUsecase := usecase.NewHealthUsecase(nil, nil, cfg)
	result := healthUsecase.GetApplicationStatus()

	assert.NotNil(t, result)
	assert.Equal(t, "test-app", result.App.Name)
	assert.Equal(t, "1.0.0", result.App.Version)
	assert.Equal(t, "test", result.App.Environment)
	assert.Equal(t, "running", result.App.Status)
	assert.NotEmpty(t, result.App.Uptime)
	assert.Equal(t, "PostgreSQL", result.Services.Database.Name)
	assert.Equal(t, domain.ServiceStatusUnhealthy, result.Services.Database.Status)

	// Test Redis service status with nil client
	assert.Equal(t, "Redis", result.Services.Redis.Name)
	assert.Equal(t, domain.ServiceStatusUnhealthy, result.Services.Redis.Status)

	assert.NotEmpty(t, result.Dependencies)
	assert.Greater(t, len(result.Dependencies), 0)

	for _, dep := range result.Dependencies {
		assert.NotEmpty(t, dep.Name)
		assert.NotEmpty(t, dep.Version)
		assert.Equal(t, domain.ServiceStatusLoaded, dep.Status)
	}
}

func TestHealthUsecase_FormatDuration(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
		},
	}

	healthUsecase := usecase.NewHealthUsecase(nil, nil, cfg)

	basicHealth := healthUsecase.GetBasicHealth()
	assert.NotNil(t, basicHealth)

	appStatus := healthUsecase.GetApplicationStatus()
	assert.NotEmpty(t, appStatus.App.Uptime)
	assert.Contains(t, []string{"s", "m", "h"}, appStatus.App.Uptime[len(appStatus.App.Uptime)-1:])
}
