import { kafka } from "./kafka.js"

const consumer = kafka.consumer({
    groupId: process.argv[2] as string
})

async function main() {
    await consumer.connect();

    console.log("consumer connected")

    await consumer.subscribe({
        topic: "orders.created",
        fromBeginning: true
    })

    await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
            console.log({
                topic,
                partition,
                key: message.key?.toString(),
                value: message.value?.toString(),
            });
        }
    })
}

main().catch(console.error);