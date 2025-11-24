import Redis from 'ioredis';
import { REDIS_URL } from './env.js';

const redisClient = new Redis(REDIS_URL);

// initialize a redis client to store active websocket connections
redisClient.on('connect', () => {
    console.log('✅ Connected to Redis server');
});

redisClient.on('error', (err) => {
    console.error('❌ Redis connection error:', err);
});

// initialize a redis pubsub subscriber to listen for messages
const redisSubscriber = new Redis(REDIS_URL);

redisSubscriber.on('connect', () => {
    console.log('✅ Redis subscriber connected to Redis server');
});

redisSubscriber.on('error', (err) => {
    console.error('❌ Redis subscriber connection error:', err);
});

export { redisClient, redisSubscriber };
