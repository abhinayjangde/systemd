interface NotificationEvent {
    eventId: string;
    eventType: string;

    actorId: string;
    recipientId: string;

    data: Record<string, unknown>;

    createdAt: Date;
}