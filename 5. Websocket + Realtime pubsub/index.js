import dotenv from "dotenv";
import http from "http";
import express from "express";
import { initWebSocketServer } from "./src/websocket/server.js";
import { PORT, REDIS_TASK_CHANNEL } from "./src/config/env.js";
import connectionRegistry from "./src/services/connectionRegistry.js";

dotenv.config();
const app = express();

const server = http.createServer(app);
initWebSocketServer(server);

// Graceful shutdown
process.on("SIGTERM", shutdown);
process.on("SIGINT", shutdown);

async function shutdown() {
    console.log("Shutting down gracefully...");

    // Close all WebSocket connections
    wss.clients.forEach((ws) => {
        ws.close();
    });

    // Close Redis connections
    await redisClient.quit();
    await redisSubscriber.quit();

    server.close(() => {
        console.log("Server closed");
        process.exit(0);
    });
}

// Start server
server.listen(PORT, () => {
    console.log(`WebSocket server running on port ${PORT}`);
    console.log(`Redis subscriber listening on channel: ${REDIS_TASK_CHANNEL}`);
});

// Health check endpoint
server.on("request", (req, res) => {
    if (req.url === "/health") {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
            JSON.stringify({
                status: "ok",
                connections: connectionRegistry.activeConnections.size,
            })
        );
    }
});
