import axios from 'axios';

// 设置API基础URL
const API_BASE_URL = 'http://localhost:8080';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 发送聊天消息的接口
export const sendChatMessage = async (message: string) => {
  try {
    const response = await apiClient.post('/chat', { 
      Messages: message 
    });
    return response.data;
  } catch (error) {
    console.error('Error sending chat message:', error);
    throw error;
  }
};

// 健康检查接口
export const checkHealth = async () => {
  try {
    const response = await apiClient.get('/health');
    return response.data;
  } catch (error) {
    console.error('Error checking health:', error);
    throw error;
  }
}; 