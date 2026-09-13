ALTER TABLE notifications
ADD COLUMN event_id UUID;

CREATE UNIQUE INDEX idx_notifications_event_id
ON notifications(event_id);