package domain

import "time"

type TicketStatus string

const (
	New        TicketStatus = "new"
	InProgress TicketStatus = "in_progress"
	Resolved   TicketStatus = "resolved"
	Closed     TicketStatus = "closed"
)

type Ticket struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	Subject     string       `json:"subject"`
	Description string       `json:"description"`
	Status      TicketStatus `json:"status"`
	Priority    int          `json:"priority"` // 1 - cao, 2 - trung bình, 3 - thấp
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
