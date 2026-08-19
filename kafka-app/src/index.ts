import { kafka } from "./lib/kafka.js";

async function init(){
    const admin = kafka.admin();
    console.log("admin connecting...");

    await admin.connect();
    console.log("admin connected!");
    
    await admin.createTopics({
        topics: [{ topic: "rider-updates", numPartitions:2 }]
    });
    console.log("Created topic!");
    await admin.disconnect();
    console.log("admin disconnected!");
}

// init().catch(console.error);