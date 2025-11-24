//  service for crud of connections in redis

import { redisClient as redis } from '../config/redis.js';

class ConnectionRegistry {
    constructor() {
        this.activeConnections = new Map();
        this.expirationTime = 60 * 60 * 24; // 24 hours
    }

    getRedisKey(authId) {
        return `ws:connections:${authId}`;
    }

    async addConnection(authId, value, ws) {
        try {
            this.activeConnections.set(authId, ws);
            await redis.set(this.getRedisKey(authId), value);
            console.log(`Connection ${authId} added`);
        } catch (error) {
            console.error(`Error adding connection ${authId}: ${error.message}`);
        }
    }

    async removeConnection(authId) {
        try {
            this.activeConnections.delete(authId);
            await redis.del(this.getRedisKey(authId));
        } catch (error) {
            console.error(`Error removing connection ${authId}: ${error.message}`);
        }
    }

    getConnection(authId) {
        try {
            console.log('Getting connection for authId:', String(authId));
            console.log(this.activeConnections.get(authId));
            return this.activeConnections.get(String(authId));
        } catch (error) {
            console.error(`Error getting connection ${authId}: ${error.message}`);
        }
    }
}

export default new ConnectionRegistry();
