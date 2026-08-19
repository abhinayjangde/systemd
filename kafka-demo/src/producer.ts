import {kafka} from './kafka.js';

const producer = kafka.producer();

async function main(){
    await producer.connect();

    console.log("producer connected")

    await producer.send({
        topic: "orders.created",
        messages: [
    {
      key: "101",
      value: JSON.stringify({
        orderId: "101",
        userId: "u1",
        amount: 500,
      }),
    },
    {
      key: "102",
      value: JSON.stringify({
        orderId: "102",
        userId: "u2",
        amount: 750,
      }),
    },
    {
      key: "103",
      value: JSON.stringify({
        orderId: "103",
        userId: "u3",
        amount: 1200,
      }),
    },
  ]
    })

    console.log("order published")

    await producer.disconnect()
}

main().catch(console.error)