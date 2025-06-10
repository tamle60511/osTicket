package domain

import "time"

type Reply struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
