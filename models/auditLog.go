package models

import (
	"time"

	"github.com/viettranx/service-context/core"
)

type AuditLog struct {
	core.SQLModel
	Time         time.Time `json:"time" gorm:"column:time;type:timestamp;not null" db:"time"`
	IPAddress    string    `json:"ip_address" gorm:"column:ip_address;type:varchar(45);not null;index" db:"ip_address"`
	APICall      string    `json:"api_call" gorm:"column:api_call;type:varchar(255);not null;index" db:"api_call"`
	Method       string    `json:"method" gorm:"column:method;type:varchar(10);not null" db:"method"`
	Status       int       `json:"status" gorm:"column:status;type:int;not null" db:"status"`
	ResponseTime int64     `json:"response_time" gorm:"column:response_time;type:bigint;not null;comment:Response time in milliseconds" db:"response_time"`
	UserID       *string   `json:"user_id" gorm:"column:user_id;type:varchar(255);index" db:"user_id"`
}

func (AuditLog) TableName() string { return "audit_logs" }

func NewAuditLog(ipAddress, apiCall, method string, status int, responseTime int64, userID *string) AuditLog {
	return AuditLog{
		SQLModel:     core.NewSQLModel(),
		Time:         time.Now(),
		IPAddress:    ipAddress,
		APICall:      apiCall,
		Method:       method,
		Status:       status,
		ResponseTime: responseTime,
		UserID:       userID,
	}
}
