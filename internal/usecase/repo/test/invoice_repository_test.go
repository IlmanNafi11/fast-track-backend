package repo_test

import (
	"fiber-boiler-plate/internal/usecase/repo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisRepository) GetJSON(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockRedisRepository) SetJSON(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisRepository) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepository) Increment(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisRepository) Decrement(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisRepository) SetExpire(key string, ttl time.Duration) error {
	args := m.Called(key, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) GetTTL(key string) (time.Duration, error) {
	args := m.Called(key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisRepository) GetKeys(pattern string) ([]string, error) {
	args := m.Called(pattern)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisRepository) FlushAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewInvoiceRepository(t *testing.T) {
	mockRedis := new(MockRedisRepository)

	invoiceRepo := repo.NewInvoiceRepository(nil, mockRedis)

	assert.NotNil(t, invoiceRepo)
}

func TestInvoiceRepository_Initialization(t *testing.T) {
	mockRedis := new(MockRedisRepository)

	// Test that repository can be initialized with nil DB for unit testing
	invoiceRepo := repo.NewInvoiceRepository(nil, mockRedis)

	assert.NotNil(t, invoiceRepo)

	// Test that repository interface is correctly implemented
	var repo repo.InvoiceRepository = invoiceRepo
	assert.NotNil(t, repo)
}

func TestInvoiceRepository_Interface(t *testing.T) {
	mockRedis := new(MockRedisRepository)
	invoiceRepo := repo.NewInvoiceRepository(nil, mockRedis)

	// Test that all interface methods exist
	assert.NotNil(t, invoiceRepo)

	// This test ensures that the repository implements the interface correctly
	var _ repo.InvoiceRepository = invoiceRepo
}
