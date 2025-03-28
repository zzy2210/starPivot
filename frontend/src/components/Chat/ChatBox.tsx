import React, { useState, useEffect, useRef } from 'react';
import { sendChatMessage } from '../../services/api';
import './ChatBox.css';

interface Message {
  id: number;
  text: string;
  isUser: boolean;
}

const ChatBox: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 自动滚动到最新消息
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!inputValue.trim()) return;
    
    // 添加用户消息
    const userMessage: Message = {
      id: Date.now(),
      text: inputValue,
      isUser: true,
    };
    
    setMessages((prevMessages) => [...prevMessages, userMessage]);
    setInputValue('');
    setIsLoading(true);
    
    try {
      // 发送消息到后端
      const response = await sendChatMessage(inputValue);
      
      // 添加AI响应
      const aiMessage: Message = {
        id: Date.now() + 1,
        text: response.message,
        isUser: false,
      };
      
      setMessages((prevMessages) => [...prevMessages, aiMessage]);
    } catch (error) {
      console.error('Failed to get response:', error);
      
      // 添加错误消息
      const errorMessage: Message = {
        id: Date.now() + 1,
        text: "抱歉，发生了错误，请稍后再试。",
        isUser: false,
      };
      
      setMessages((prevMessages) => [...prevMessages, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h2>StarPivot 聊天</h2>
      </div>
      
      <div className="chat-messages">
        {messages.length === 0 ? (
          <div className="empty-state">
            <p>发送一条消息开始聊天！</p>
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message.id}
              className={`message ${message.isUser ? 'user-message' : 'ai-message'}`}
            >
              <div className="message-content">
                <p>{message.text}</p>
              </div>
            </div>
          ))
        )}
        {isLoading && (
          <div className="message ai-message">
            <div className="message-content loading">
              <p>正在思考...</p>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>
      
      <form className="chat-input-form" onSubmit={handleSubmit}>
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          placeholder="输入消息..."
          disabled={isLoading}
        />
        <button type="submit" disabled={isLoading || !inputValue.trim()}>
          发送
        </button>
      </form>
    </div>
  );
};

export default ChatBox; 