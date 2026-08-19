import { Kafka } from "kafkajs"

const kafka = new Kafka({
  clientId: "notification-service",
  brokers: ["localhost:9092"],
});

const consumer = kafka.consumer({
  groupId: "notification-service-group"
});

async function consumeOrders() {
  await consumer.connect();

  await consumer.subscribe({
    topic: "orders",
    fromBeginning: true,
  });

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      const order = JSON.parse(message.value.toString());

      console.log({
        topic,
        partition,
        offset: message.offset,
        order,
      });
    },
  });
}

consumeOrders().catch(console.error);