import React, { useState, useEffect } from 'react';
import { getChatIds, createNewChat } from '../../services/api';
import './ChatList.css';

interface ChatListProps {
  onSelectChat: (chatId: string) => void;
  selectedChatId: string | null;
}

const ChatList: React.FC<ChatListProps> = ({ onSelectChat, selectedChatId }) => {
  const [chatIds, setChatIds] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // 加载聊天列表
  const loadChatList = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const ids = await getChatIds();
      setChatIds(ids);
      
      // 如果有聊天ID但还没有选择的聊天，则选择第一个
      if (ids.length > 0 && !selectedChatId) {
        onSelectChat(ids[0]);
      }
    } catch (err) {
      console.error('Failed to load chat list:', err);
      setError('加载聊天列表失败');
    } finally {
      setIsLoading(false);
    }
  };

  // 创建新聊天
  const handleCreateNewChat = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const { chatId } = await createNewChat();
      setChatIds(prevIds => [chatId, ...prevIds]);
      onSelectChat(chatId);
    } catch (err) {
      console.error('Failed to create new chat:', err);
      setError('创建新聊天失败');
    } finally {
      setIsLoading(false);
    }
  };

  // 组件加载时获取聊天列表
  useEffect(() => {
    loadChatList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="chat-list-container">
      <div className="chat-list-header">
        <h3>我的聊天</h3>
        <button
          className="new-chat-button"
          onClick={handleCreateNewChat}
          disabled={isLoading}
        >
          新聊天
        </button>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="chat-list">
        {isLoading && chatIds.length === 0 ? (
          <div className="loading-indicator">加载中...</div>
        ) : chatIds.length === 0 ? (
          <div className="empty-chats">
            <p>没有聊天记录</p>
            <p>点击"新聊天"开始对话</p>
          </div>
        ) : (
          <ul className="chat-items">
            {chatIds.map(id => (
              <li
                key={id}
                className={`chat-item ${selectedChatId === id ? 'selected' : ''}`}
                onClick={() => onSelectChat(id)}
              >
                <div className="chat-item-title">
                  聊天 {id.substring(0, 8)}...
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default ChatList; 