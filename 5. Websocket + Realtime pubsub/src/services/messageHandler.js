import { WebSocket } from "ws";
import connectionRegistry from "./connectionRegistry.js";

class MessageHandler {
    async sendMessage(authId, message) {
        const ws = connectionRegistry.getConnection(authId);
        console.log({
            authId,
            ws,
            a: ws?.readyState,
            b: WebSocket.OPEN
        })
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(message));
            return true;
        } else {
            console.error(`Connection for authId ${authId} is not open.`);
            return false;
        }
    }

    async broadcastMessage(message) {
        for (const [authId, ws] of connectionRegistry.activeConnections.entries()) {
            if (ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify(message));
            } else {
                console.error(`Connection for authId ${authId} is not open.`);
            }
        }
    }

    handleMessage(authId, message) {
        if (['broadcast'].includes(message.type)) {
            return this.broadcastMessage(message);
        } else {
            return this.sendMessage(authId, message);
        }
    }
}

export default new MessageHandler();