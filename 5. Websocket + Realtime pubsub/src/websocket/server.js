import { WebSocketServer } from 'ws';
import startTaskSubscriber from '../services/taskSubscriber.js';
import handleConnection from './handler.js';
import MessageHandler from '../services/messageHandler.js';

export const initWebSocketServer = (server) => {
    const wss = new WebSocketServer({ server });

    wss.on('connection', (ws, req) => {
        handleConnection(ws, req);
    });


    startTaskSubscriber();
}