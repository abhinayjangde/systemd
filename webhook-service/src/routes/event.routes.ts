import { Router } from "express";
import { prisma } from "../lib/prisma.js";
import { deliverEvent } from "../services/webhook.service.js";

const router = Router();

router.post("/", async (req, res) => {
    try {
        const { type, payload } = req.body;

        if (!type || !payload) {
            return res.status(400).json({
                message: "type and payload are required",
            });
        }

        const event = await prisma.event.create({
            data: {
                type,
                payload,
            },
        });

        await deliverEvent(event.id);
        
        return res.status(201).json(event);

    } catch (error) {
        console.error(error);

        return res.status(500).json({
            message: "Internal server error",
        });
    }
});

export default router;