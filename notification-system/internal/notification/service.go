package notification

import (
	"context"

	"github.com/abhinayjangde/notification-system/internal/events"
)

type Event = events.NotificationEvent

type Repository interface {
	CreateFromEvent(ctx context.Context, event Event) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) ProcessEvent(
	ctx context.Context,
	event Event,
) error {
	if event.EventID == "" {
		return events.NewPermanentError("event_id is required")
	}

	if event.EventType == "" {
		return events.NewPermanentError("event_type is required")
	}

	if event.RecipientID == "" {
		return events.NewPermanentError("recipient_id is required")
	}

	if event.Title == "" {
		return events.NewPermanentError("title is required")
	}

	if event.Body == "" {
		return events.NewPermanentError("body is required")
	}

	return s.repository.CreateFromEvent(ctx, event)
}
