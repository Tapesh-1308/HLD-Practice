export const serverId = process.env.SERVER_ID || 'server-1';
export const REDIS_URL = process.env.REDIS_URL || "redis://127.0.0.1:6379";
export const JWT_SECRET = process.env.JWT_SECRET || "your_jwt_secret_key";
export const PORT = process.env.PORT || 4000;
export const REDIS_TASK_CHANNEL = process.env.REDIS_TASK_CHANNEL || "tasks_channel";