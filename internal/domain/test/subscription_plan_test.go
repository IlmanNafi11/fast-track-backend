package domain

import (
	"fiber-boiler-plate/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionPlanListRequest_SetDefaults(t *testing.T) {
	t.Run("should set default values when empty", func(t *testing.T) {
		req := &domain.SubscriptionPlanListRequest{}
		req.SetDefaults()

		assert.Equal(t, "nama", req.SortBy)
		assert.Equal(t, "asc", req.SortDirection)
		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.PerPage)
	})

	t.Run("should not override existing values", func(t *testing.T) {
		req := &domain.SubscriptionPlanListRequest{
			SortBy:        "harga",
			SortDirection: "desc",
			Page:          2,
			PerPage:       20,
		}
		req.SetDefaults()

		assert.Equal(t, "harga", req.SortBy)
		assert.Equal(t, "desc", req.SortDirection)
		assert.Equal(t, 2, req.Page)
		assert.Equal(t, 20, req.PerPage)
	})

	t.Run("should set default values for invalid page and per_page", func(t *testing.T) {
		req := &domain.SubscriptionPlanListRequest{
			Page:    0,
			PerPage: 0,
		}
		req.SetDefaults()

		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.PerPage)
	})
}

func TestSubscriptionPlanListRequest_GetOffset(t *testing.T) {
	testCases := []struct {
		name     string
		page     int
		perPage  int
		expected int
	}{
		{
			name:     "page 1 with 10 per page",
			page:     1,
			perPage:  10,
			expected: 0,
		},
		{
			name:     "page 2 with 10 per page",
			page:     2,
			perPage:  10,
			expected: 10,
		},
		{
			name:     "page 3 with 20 per page",
			page:     3,
			perPage:  20,
			expected: 40,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &domain.SubscriptionPlanListRequest{
				Page:    tc.page,
				PerPage: tc.perPage,
			}

			offset := req.GetOffset()
			assert.Equal(t, tc.expected, offset)
		})
	}
}
