import { Router } from "express";
import { prisma } from "../lib/prisma.js";

const router = Router();

router.post("/", async (req, res) => {
    try {
        const { url, secret } = req.body;

        if (!url || !secret) {
            return res.status(400).json({
                message: "url and secret are required",
            });
        }

        const webhook = await prisma.webhookEndpoint.create({
            data: {
                url,
                secret,
            },
        });

        return res.status(201).json(webhook);
    } catch (error) {
        console.error(error);

        return res.status(500).json({
            message: "Internal server error",
        });
    }
});

router.get("/", async (req, res) => {
    try {
        const webhooks = await prisma.webhookEndpoint.findMany({
            select: {
                id: true,
                url: true,
                isActive: true,
                createdAt: true
            }
        });
        return res.status(200).json(webhooks)
    } catch (error) {
        console.error(error);

        return res.status(500).json({
            message: "Internal server error",
        });
    }

})
router.get("/:id", async (req, res) => {
    try {
        const id = req.params.id;

        const webhooks = await prisma.webhookEndpoint.findUnique({
            where: {
                id
            },
            select: {
                id: true,
                url: true,
                isActive: true,
                createdAt: true
            }
        });

        if(!webhooks){
            return res.status(404).json({message:"webhook doesn't exists"})
        }
        return res.status(200).json(webhooks)
    } catch (error) {
        console.error(error);

        return res.status(500).json({
            message: "Internal server error",
        });
    }

})

export default router;