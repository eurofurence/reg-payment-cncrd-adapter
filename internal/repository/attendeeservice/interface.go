package attendeeservice

import (
	"context"
	"errors"
)

type AttendeeService interface {
	GetAttendee(ctx context.Context, id uint) (AttendeeDto, error)
}

var (
	NotFoundError   = errors.New("attendee id not found")
	DownstreamError = errors.New("downstream unavailable - see log for details")
)

// we only list the fields we may actually use

type AttendeeDto struct {
	Id                   uint   `json:"id"`       // badge number
	Nickname             string `json:"nickname"` // fan name
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Street               string `json:"street"`
	Zip                  string `json:"zip"`
	City                 string `json:"city"`
	Country              string `json:"country"` // 2 letter ISO-3166-1 country code for the address (Alpha-2 code)
	Email                string `json:"email"`
	RegistrationLanguage string `json:"registration_language"` // one out of configurable subset of RFC 5646 locales (default en-US)
}
