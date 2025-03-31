import React, { useState, useEffect } from 'react';
import { getChatIds, createNewChat, deleteChat } from '../../services/api';
import './ChatList.css';

interface ChatListProps {
  onSelectChat: (chatId: string) => void;
  selectedChatId: string | null;
  onChatDeleted?: () => void;
}

const ChatList: React.FC<ChatListProps> = ({ onSelectChat, selectedChatId, onChatDeleted }) => {
  const [chatIds, setChatIds] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState<boolean>(false);

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

  // 添加新聊天到列表（由ChatBox在新建聊天时调用）
  const addChatToList = (chatId: string) => {
    if (!chatIds.includes(chatId)) {
      setChatIds(prevIds => [chatId, ...prevIds]);
    }
  };

  // 删除聊天
  const handleDeleteChat = (e: React.MouseEvent, chatId: string) => {
    e.stopPropagation(); // 防止点击事件冒泡到聊天项
    setShowDeleteConfirm(chatId);
  };

  // 确认删除聊天
  const confirmDeleteChat = async (chatId: string) => {
    setIsDeleting(true);
    setError(null);
    
    try {
      const success = await deleteChat(chatId);
      if (success) {
        // 从列表中移除
        setChatIds(prevIds => prevIds.filter(id => id !== chatId));
        
        // 如果删除的是当前选中的聊天，选择另一个
        if (selectedChatId === chatId) {
          const remainingIds = chatIds.filter(id => id !== chatId);
          if (remainingIds.length > 0) {
            onSelectChat(remainingIds[0]);
          } else {
            onSelectChat('');
            if (onChatDeleted) {
              onChatDeleted();
            }
          }
        }
      } else {
        setError('删除聊天失败');
      }
    } catch (err) {
      console.error('Failed to delete chat:', err);
      setError('删除聊天失败');
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(null);
    }
  };

  // 取消删除确认
  const cancelDeleteChat = (e: React.MouseEvent) => {
    e.stopPropagation(); // 防止点击事件冒泡
    setShowDeleteConfirm(null);
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
          disabled={isLoading || isDeleting}
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
                
                <button 
                  className="chat-delete-btn"
                  onClick={(e) => handleDeleteChat(e, id)}
                  disabled={isDeleting}
                  title="删除聊天"
                >
                  ×
                </button>
                
                {showDeleteConfirm === id && (
                  <div className="delete-confirm-popup">
                    <p>确定要删除这个聊天吗？</p>
                    <div className="delete-confirm-buttons">
                      <button 
                        onClick={() => confirmDeleteChat(id)}
                        disabled={isDeleting}
                      >
                        {isDeleting ? '删除中...' : '确定'}
                      </button>
                      <button 
                        onClick={cancelDeleteChat}
                        disabled={isDeleting}
                      >
                        取消
                      </button>
                    </div>
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default ChatList; 