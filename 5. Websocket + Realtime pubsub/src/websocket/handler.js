import ConnectionRegistry from '../services/connectionRegistry.js';

const handleConnection = async (ws, req) => {
    const url = new URL(req.url, `http://${req.headers.host}`);
    const authId = url.searchParams.get('authId');

    await ConnectionRegistry.addConnection(authId, {}, ws);

    ws.on('message', (message) => {
        try {
            const data = JSON.parse(message);
            if (data.type === 'ping') {
                ws.send(JSON.stringify({ type: 'pong' }));
            }
        } catch (error) {
            console.error(`Error processing message from ${authId}: ${error.message}`);
        }
    });

    ws.on('close', async () => {
        await ConnectionRegistry.removeConnection(authId, ws);
        console.log(`Connection ${authId} closed`);
    });

    ws.on('error', (error) => {
        console.error(`Connection ${authId} error: ${error.message}`);
    });
};

export default handleConnection;