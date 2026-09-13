package notification

import (
	"context"
	"fmt"

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
		return fmt.Errorf("event_id is required")
	}

	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}

	if event.RecipientID == "" {
		return fmt.Errorf("recipient_id is required")
	}

	return s.repository.CreateFromEvent(ctx, event)
}
