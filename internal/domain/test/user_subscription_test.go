package domain_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserSubscription_BeforeCreate(t *testing.T) {
	subscription := &domain.UserSubscription{}

	err := subscription.BeforeCreate(nil)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, subscription.ID)
}

func TestUserSubscription_BeforeCreate_ExistingID(t *testing.T) {
	existingID := uuid.New()
	subscription := &domain.UserSubscription{ID: existingID}

	err := subscription.BeforeCreate(nil)

	assert.NoError(t, err)
	assert.Equal(t, existingID, subscription.ID)
}

func TestUserSubscriptionListRequest_SetDefaults(t *testing.T) {
	req := &domain.UserSubscriptionListRequest{}

	req.SetDefaults()

	assert.Equal(t, "nama", req.SortBy)
	assert.Equal(t, "asc", req.SortDirection)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PerPage)
}

func TestUserSubscriptionListRequest_GetOffset(t *testing.T) {
	req := &domain.UserSubscriptionListRequest{
		Page:    3,
		PerPage: 15,
	}

	offset := req.GetOffset()

	assert.Equal(t, 30, offset) // (3-1) * 15 = 30
}

func TestUpdateUserSubscriptionRequest_Validation(t *testing.T) {
	req := &domain.UpdateUserSubscriptionRequest{
		Action: "pause",
	}

	assert.Equal(t, "pause", req.Action)
}

func TestUpdatePaymentMethodRequest_Validation(t *testing.T) {
	req := &domain.UpdatePaymentMethodRequest{
		PaymentMethod: "Credit Card",
	}

	assert.Equal(t, "Credit Card", req.PaymentMethod)
}
