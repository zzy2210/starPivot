import axios from 'axios';

// 定义接口类型
export interface Message {
  id: number;
  text: string;
  isUser: boolean;
  timestamp: number;
}

export interface Chat {
  id: string;
  messages: Message[];
}

// 设置API基础URL
const API_BASE_URL = 'http://localhost:8080';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 获取所有聊天ID列表
export const getChatIds = async (): Promise<string[]> => {
  try {
    const response = await apiClient.get('/chat/chat/ids', {
      params: {
        Username: localStorage.getItem('username') || '匿名用户'
      }
    });
    // 后端返回 ChatIDs 而不是 chatIds
    if (response.data && Array.isArray(response.data.ChatIDs)) {
      return response.data.ChatIDs;
    } else {
      console.error('Unexpected response format:', response.data);
      return [];
    }
  } catch (error) {
    console.error('Error getting chat IDs:', error);
    return [];
  }
};

// 获取特定ID的聊天内容
export const getChatById = async (chatId: string): Promise<Chat> => {
  try {
    const response = await apiClient.get(`/chat/chat/${chatId}`, {
      params: {
        Username: localStorage.getItem('username') || '匿名用户',
        ChatID: chatId
      }
    });
    
    // 转换后端的消息格式为前端格式
    const messages = Array.isArray(response.data.Messages) 
      ? response.data.Messages.map((msg: any, index: number) => ({
          id: index,
          text: msg.Content || msg.content || '',
          isUser: msg.Role === 'user' || msg.role === 'user',
          timestamp: new Date().getTime()
        }))
      : [];
      
    return {
      id: chatId,
      messages: messages
    };
  } catch (error) {
    console.error(`Error getting chat ${chatId}:`, error);
    throw error;
  }
};

// 创建新的聊天
export const createNewChat = async (): Promise<{chatId: string}> => {
  try {
    const response = await apiClient.post('/chat/chat/new', {
      Username: localStorage.getItem('username') || '匿名用户'
    });
    // 后端返回 ChatID 而不是 chatId
    return { chatId: response.data.ChatID };
  } catch (error) {
    console.error('Error creating new chat:', error);
    throw error;
  }
};

// 发送聊天消息的接口
export const sendChatMessage = async (message: string, chatId?: string) => {
  try {
    const endpoint = '/chat/chat';
    const response = await apiClient.post(endpoint, { 
      Messages: message,
      ChatID: chatId || ''
    });
    return {
      message: response.data.Messages
    };
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