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
  Messages: Message[];
}

// 设置API基础URL
const API_BASE_URL = 'http://localhost:8080';

// 获取用户名
const getUsername = () => localStorage.getItem('username') || '匿名用户';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器添加用户名请求头
apiClient.interceptors.request.use(config => {
  config.headers['X-Username'] = getUsername();
  return config;
});

// 获取所有聊天ID列表
export const getChatIds = async (): Promise<string[]> => {
  try {
    const response = await apiClient.get('/chat/ids');
    if (response.data && response.data.code === 200 && Array.isArray(response.data.data.chatIDs)) {
      return response.data.data.chatIDs;
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
    const response = await apiClient.get(`/chat/${chatId}`);
    
    if (response.data && response.data.code === 200) {
      // 转换后端的消息格式为前端格式
      const messages = Array.isArray(response.data.data.Messages)
        ? response.data.data.Messages.map((msg: any, index: number) => ({
            id: index,
            text: msg.content || '',
            isUser: msg.role === 'user',
            timestamp: new Date().getTime()
          }))
        : [];
        
      return {
        id: chatId,
        Messages: messages
      };
    } else {
      console.error('Unexpected response format:', response.data);
      return { id: chatId, Messages: [] };
    }
  } catch (error) {
    console.error(`Error getting chat ${chatId}:`, error);
    throw error;
  }
};

// 创建新的聊天
export const createNewChat = async (): Promise<{chatId: string}> => {
  try {
    const response = await apiClient.post('/chat/new');
    if (response.data && response.data.code === 200) {
      return { chatId: response.data.data.chatID };
    } else {
      console.error('Unexpected response format:', response.data);
      throw new Error('Failed to create new chat');
    }
  } catch (error) {
    console.error('Error creating new chat:', error);
    throw error;
  }
};

// 发送聊天消息的接口
export const sendChatMessage = async (message: string, chatId?: string) => {
  try {
    const response = await apiClient.post('/chat/chat', { 
      messages: message,
      chatID: chatId || undefined
    });
    
    if (response.data && response.data.code === 200) {
      return {
        message: response.data.data.messages,
        chatId: response.data.data.chatID
      };
    } else {
      console.error('Unexpected response format:', response.data);
      throw new Error('Failed to send chat message');
    }
  } catch (error) {
    console.error('Error sending chat message:', error);
    throw error;
  }
};

// 删除聊天
export const deleteChat = async (chatId: string): Promise<boolean> => {
  try {
    const response = await apiClient.delete(`/chat/${chatId}`);
    return response.data && response.data.code === 200;
  } catch (error) {
    console.error(`Error deleting chat ${chatId}:`, error);
    return false;
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