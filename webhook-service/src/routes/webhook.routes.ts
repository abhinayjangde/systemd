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

export default router;