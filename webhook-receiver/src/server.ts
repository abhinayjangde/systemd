import express from "express";

const app = express();

app.use(express.json());

app.post("/webhook", (req, res) => {
    console.log("Webhook received!");

    console.log(req.body);

    return res.status(200).json({
        received: true,
    });
});

app.listen(4000, () => {
    console.log("Receiver running on http://localhost:4000");
});