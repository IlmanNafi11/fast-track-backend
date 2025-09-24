package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSubscription struct {
	ID                 uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	UserID             uint             `json:"user_id" gorm:"not null;index"`
	User               User             `json:"user" gorm:"foreignKey:UserID"`
	SubscriptionPlanID uuid.UUID        `json:"subscription_plan_id" gorm:"type:uuid;not null;index"`
	SubscriptionPlan   SubscriptionPlan `json:"subscription_plan" gorm:"foreignKey:SubscriptionPlanID"`
	Status             string           `json:"status" gorm:"not null;default:'trialing';check:status IN ('active','paused','trialing','canceled','ended')"`
	CurrentPeriodStart time.Time        `json:"current_period_start" gorm:"not null"`
	CurrentPeriodEnd   time.Time        `json:"current_period_end" gorm:"not null"`
	TrialEnd           *time.Time       `json:"trial_end" gorm:"index"`
	PaymentMethod      string           `json:"payment_method" gorm:"not null"`
	PaymentStatus      string           `json:"payment_status" gorm:"not null;default:'pending';check:payment_status IN ('paid','pending','failed','refunded')"`
	CancelAtPeriodEnd  bool             `json:"cancel_at_period_end" gorm:"default:false"`
	CanceledAt         *time.Time       `json:"canceled_at"`
	EndedAt            *time.Time       `json:"ended_at"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

func (us *UserSubscription) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

type UserSubscriptionListRequest struct {
	Search        *string `json:"search" query:"search"`
	Status        *string `json:"status" query:"status" validate:"omitempty,oneof=active paused trialing canceled ended"`
	PaymentMethod *string `json:"payment_method" query:"payment_method"`
	SortBy        string  `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=nama email status created_at"`
	SortDirection string  `json:"sort_direction" query:"sort_direction" validate:"omitempty,oneof=asc desc"`
	Page          int     `json:"page" query:"page" validate:"min=1"`
	PerPage       int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func (r *UserSubscriptionListRequest) GetOffset() int {
	return (r.Page - 1) * r.PerPage
}

func (req *UserSubscriptionListRequest) SetDefaults() {
	if req.SortBy == "" {
		req.SortBy = "nama"
	}
	if req.SortDirection == "" {
		req.SortDirection = "asc"
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 10
	}
}

type UpdateUserSubscriptionRequest struct {
	Action string  `json:"action" validate:"required,oneof=pause activate cancel"`
	Reason *string `json:"reason" validate:"omitempty,max=500"`
}

type UpdatePaymentMethodRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,min=1,max=100"`
}

type UserSubscriptionResponse struct {
	ID                 uuid.UUID            `json:"id"`
	User               UserInfo             `json:"user"`
	SubscriptionPlan   SubscriptionPlanInfo `json:"subscription_plan"`
	Status             string               `json:"status"`
	CurrentPeriodStart time.Time            `json:"current_period_start"`
	CurrentPeriodEnd   time.Time            `json:"current_period_end"`
	PaymentMethod      string               `json:"payment_method"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

type UserSubscriptionDetailResponse struct {
	ID                 uuid.UUID              `json:"id"`
	User               UserInfoDetail         `json:"user"`
	SubscriptionPlan   SubscriptionPlanDetail `json:"subscription_plan"`
	Status             string                 `json:"status"`
	CurrentPeriodStart time.Time              `json:"current_period_start"`
	CurrentPeriodEnd   time.Time              `json:"current_period_end"`
	TrialEnd           *time.Time             `json:"trial_end"`
	PaymentMethod      string                 `json:"payment_method"`
	PaymentStatus      string                 `json:"payment_status"`
	CancelAtPeriodEnd  bool                   `json:"cancel_at_period_end"`
	CanceledAt         *time.Time             `json:"canceled_at"`
	EndedAt            *time.Time             `json:"ended_at"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Nama  string `json:"nama"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserInfoDetail struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type SubscriptionPlanInfo struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Harga float64   `json:"harga"`
}

type SubscriptionPlanDetail struct {
	ID            uuid.UUID `json:"id"`
	Kode          string    `json:"kode"`
	Nama          string    `json:"nama"`
	Harga         float64   `json:"harga"`
	Interval      string    `json:"interval"`
	HariPercobaan int       `json:"hari_percobaan"`
	Status        string    `json:"status"`
}

type UserSubscriptionStatistics struct {
	TotalSubscriptions    int64                    `json:"total_subscriptions"`
	ActiveSubscriptions   int64                    `json:"active_subscriptions"`
	PausedSubscriptions   int64                    `json:"paused_subscriptions"`
	TrialingSubscriptions int64                    `json:"trialing_subscriptions"`
	PaymentMethods        []PaymentMethodStatistic `json:"payment_methods"`
	MonthlyRevenue        float64                  `json:"monthly_revenue"`
	YearlyRevenue         float64                  `json:"yearly_revenue"`
}

type PaymentMethodStatistic struct {
	Method     string  `json:"method"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}
