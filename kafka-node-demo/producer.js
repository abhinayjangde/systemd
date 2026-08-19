import { Kafka, Partitioners } from "kafkajs"

const kafka = new Kafka({
    clientId: "order-service",
    brokers: ["localhost:9092"]
})

const producer = kafka.producer({
    createPartitioner: Partitioners.DefaultPartitioner
})

async function sendOrder() {
    await producer.connect()
    const result = await producer.send({
        topic: "orders",
        messages: [
            {
                key: "user-43",
                value: JSON.stringify({
                    eventType: "OrderCreated",
                    orderId: "ORD-106",
                    userId: "user-43",
                    amount: 999,
                }),
            },
        ],
    })

    console.log("Message sent:", result);
    await producer.disconnect();
}

sendOrder().catch(console.error);