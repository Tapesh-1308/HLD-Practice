import { REDIS_TASK_CHANNEL } from "../config/env.js";
import { redisSubscriber } from "../config/redis.js";
import messageHandler from "./messageHandler.js";

const startTaskSubscriber = () => {
    redisSubscriber.subscribe(REDIS_TASK_CHANNEL, (err, count) => {
        if (err) {
            console.error('❌ Failed to subscribe: ', err);
            return;
        }
        console.log(`✅ Subscribed successfully! This client is currently subscribed to ${count} channels.`);
    });

    redisSubscriber.on('message', async (channel, message) => {
        if (channel === REDIS_TASK_CHANNEL) {
            try {
                const {authId, data} = JSON.parse(message);
                console.log(`📩 Received message for authId ${authId}: `, data);
                const send = await messageHandler.handleMessage(authId, data);
                console.log(messageHandler.handleMessage, send)
                if (send) {
                    console.log('✅ Message sent successfully');
                } else {
                    console.error('❌ Failed to send message');
                }
            } catch (error) {
                console.error('❌ Failed to parse message: ', error);
            }
        }
    });
};

export default startTaskSubscriber;