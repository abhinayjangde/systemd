import { kafka } from "./lib/kafka.js";

async function init(){
    const producer = kafka.producer();
    console.log("producer connecting...");

    await producer.connect();
    console.log("producer connected!");
    
    await producer.send({
        topic: "rider-updates",
        messages: [
            { value: "Hello KafkaJS user!" },
        ],
    });
    console.log("Message sent!");
    await producer.disconnect();
    console.log("producer disconnected!");
}

init().catch(console.error);