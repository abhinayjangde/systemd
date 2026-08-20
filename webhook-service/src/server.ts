import express from "express";
import dotenv from "dotenv";
import webhookRoutes from "./routes/webhook.routes.js";
import eventRoutes from "./routes/event.routes.js";

dotenv.config();

const app = express();

app.use(express.json());

app.get("/", (req, res) => {
  res.json({
    message: "Webhook service is running",
  });
});

app.use("/webhooks", webhookRoutes);
app.use("/events", eventRoutes);
const PORT = process.env.PORT || 3000;

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});