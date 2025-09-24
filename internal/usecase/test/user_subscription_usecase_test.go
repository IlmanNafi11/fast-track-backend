package usecase

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserSubscriptionRepository struct {
	mock.Mock
}

func (m *MockUserSubscriptionRepository) GetAll(req *domain.UserSubscriptionListRequest) ([]*domain.UserSubscription, int, error) {
	args := m.Called(req)
	return args.Get(0).([]*domain.UserSubscription), args.Int(1), args.Error(2)
}

func (m *MockUserSubscriptionRepository) GetByID(id string) (*domain.UserSubscription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserSubscription), args.Error(1)
}

func (m *MockUserSubscriptionRepository) UpdateStatus(id string, status string, reason *string) error {
	args := m.Called(id, status, reason)
	return args.Error(0)
}

func (m *MockUserSubscriptionRepository) UpdatePaymentMethod(id string, paymentMethod string) error {
	args := m.Called(id, paymentMethod)
	return args.Error(0)
}

func (m *MockUserSubscriptionRepository) GetStatistics() (*domain.UserSubscriptionStatistics, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserSubscriptionStatistics), args.Error(1)
}

func (m *MockUserSubscriptionRepository) Create(subscription *domain.UserSubscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockUserSubscriptionRepository) Update(subscription *domain.UserSubscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockUserSubscriptionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserSubscriptionUsecase_GetAll_Success(t *testing.T) {
	mockRepo := new(MockUserSubscriptionRepository)

	useCase := usecase.NewUserSubscriptionUsecase(mockRepo, nil, nil)

	now := time.Now()
	subscriptions := []*domain.UserSubscription{
		{
			ID:                 uuid.New(),
			UserID:             1,
			User:               domain.User{ID: 1, Name: "John Doe", Email: "john@example.com", IsActive: true},
			SubscriptionPlanID: uuid.New(),
			SubscriptionPlan:   domain.SubscriptionPlan{ID: uuid.New(), Nama: "PRO Monthly", Harga: 99000},
			Status:             "active",
			CurrentPeriodStart: now,
			CurrentPeriodEnd:   now.AddDate(0, 1, 0),
			PaymentMethod:      "Bank Transfer",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	req := &domain.UserSubscriptionListRequest{
		Page:          1,
		PerPage:       10,
		SortBy:        "nama",
		SortDirection: "asc",
	}

	mockRepo.On("GetAll", req).Return(subscriptions, 1, nil)

	result, meta, err := useCase.GetAll(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, meta)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, meta.CurrentPage)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 1, meta.TotalRecords)
	assert.Equal(t, 10, meta.PerPage)

	mockRepo.AssertExpectations(t)
}

func TestUserSubscriptionUsecase_GetByID_Success(t *testing.T) {
	mockRepo := new(MockUserSubscriptionRepository)

	useCase := usecase.NewUserSubscriptionUsecase(mockRepo, nil, nil)

	subscriptionID := uuid.New().String()
	subscription := &domain.UserSubscription{
		ID:                 uuid.MustParse(subscriptionID),
		UserID:             1,
		User:               domain.User{ID: 1, Name: "John Doe", Email: "john@example.com", IsActive: true},
		SubscriptionPlanID: uuid.New(),
		Status:             "active",
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
		PaymentMethod:      "Bank Transfer",
	}

	mockRepo.On("GetByID", subscriptionID).Return(subscription, nil)

	result, err := useCase.GetByID(subscriptionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, subscriptionID, result.ID.String())
	assert.Equal(t, "John Doe", result.User.Nama)

	mockRepo.AssertExpectations(t)
}

func TestUserSubscriptionUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserSubscriptionRepository)

	useCase := usecase.NewUserSubscriptionUsecase(mockRepo, nil, nil)

	subscriptionID := uuid.New().String()
	mockRepo.On("GetByID", subscriptionID).Return(nil, gorm.ErrRecordNotFound)

	result, err := useCase.GetByID(subscriptionID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "subscription pengguna tidak ditemukan", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserSubscriptionUsecase_GetStatistics_Success(t *testing.T) {
	mockRepo := new(MockUserSubscriptionRepository)

	useCase := usecase.NewUserSubscriptionUsecase(mockRepo, nil, nil)

	stats := &domain.UserSubscriptionStatistics{
		TotalSubscriptions:    150,
		ActiveSubscriptions:   120,
		PausedSubscriptions:   20,
		TrialingSubscriptions: 10,
		PaymentMethods: []domain.PaymentMethodStatistic{
			{Method: "Bank Transfer", Count: 75, Percentage: 50.0},
		},
		MonthlyRevenue: 9500000,
		YearlyRevenue:  114000000,
	}

	mockRepo.On("GetStatistics").Return(stats, nil)

	result, err := useCase.GetStatistics()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(150), result.TotalSubscriptions)
	assert.Equal(t, int64(120), result.ActiveSubscriptions)
	assert.Equal(t, float64(9500000), result.MonthlyRevenue)

	mockRepo.AssertExpectations(t)
}
