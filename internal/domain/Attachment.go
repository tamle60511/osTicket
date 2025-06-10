package domain

import "time"

type Attachment struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
