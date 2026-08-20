import axios from "axios";
import { prisma } from "../lib/prisma.js";

export async function deliverEvent(eventId: string) {
    const event = await prisma.event.findUnique({
        where: {
            id: eventId,
        },
    });

    if (!event) {
        throw new Error("Event not found");
    }

    const webhooks = await prisma.webhookEndpoint.findMany({
        where: {
            isActive: true,
        },
    });

    for (const webhook of webhooks) {
        const delivery = await prisma.webhookDelivery.create({
            data: {
                eventId: event.id,
                webhookId: webhook.id,
            },
        });

        try {
            await axios.post(
                webhook.url,
                {
                    id: event.id,
                    type: event.type,
                    payload: event.payload,
                },
                {
                    timeout: 5000,
                }
            );

            await prisma.webhookDelivery.update({
                where: {
                    id: delivery.id,
                },
                data: {
                    status: "DELIVERED",
                    attempts: 1,
                    deliveredAt: new Date(),
                },
            });

        } catch (error) {
            await prisma.webhookDelivery.update({
                where: {
                    id: delivery.id,
                },
                data: {
                    status: "FAILED",
                    attempts: 1,
                    lastError: error instanceof Error
                        ? error.message
                        : "Unknown error",
                },
            });
        }
    }
}