package service

import (
	"time"

	"api.mijkomp.com/models/response"
)

type DashboardService interface {
	GetSummary(currentUserId uint, fromDate, toDate time.Time) response.Dashboard
}
