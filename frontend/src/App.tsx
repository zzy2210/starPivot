import React, { useState, useEffect } from 'react';
import Header from './components/Layout/Header';
import Sidebar from './components/Layout/Sidebar';
import ChatBox from './components/Chat/ChatBox';
import ChatList from './components/Chat/ChatList';
import { checkHealth } from './services/api';
import './App.css';

function App() {
  const [isServerConnected, setIsServerConnected] = useState<boolean | null>(null);
  const [username, setUsername] = useState<string>('');
  const [selectedChatId, setSelectedChatId] = useState<string | null>(null);
  
  // 处理用户名变更
  const handleUsernameChange = (newUsername: string) => {
    setUsername(newUsername);
    localStorage.setItem('username', newUsername);
  };

  // 处理聊天选择
  const handleSelectChat = (chatId: string) => {
    setSelectedChatId(chatId);
  };
  
  // 处理新聊天创建（从ChatBox组件调用）
  const handleChatCreated = (chatId: string) => {
    setSelectedChatId(chatId);
  };

  // 处理聊天删除
  const handleChatDeleted = () => {
    setSelectedChatId(null);
  };

  // 初始化用户名
  useEffect(() => {
    const storedUsername = localStorage.getItem('username');
    if (storedUsername) {
      setUsername(storedUsername);
    } else {
      setUsername('匿名用户');
      localStorage.setItem('username', '匿名用户');
    }
  }, []);

  // 检查服务器健康状态
  useEffect(() => {
    const checkServerHealth = async () => {
      try {
        const healthStatus = await checkHealth();
        setIsServerConnected(healthStatus.status === 'ok');
      } catch (error) {
        console.error('Health check failed:', error);
        setIsServerConnected(false);
      }
    };

    checkServerHealth();
    
    // 每30秒检查一次服务器连接状态
    const intervalId = setInterval(checkServerHealth, 30000);
    
    return () => clearInterval(intervalId);
  }, []);

  return (
    <div className="App">
      <Header onUsernameChange={handleUsernameChange} username={username} />
      
      <div className="server-status-container">
        {isServerConnected === null ? (
          <div className="server-status checking">正在连接服务器...</div>
        ) : isServerConnected ? (
          <div className="server-status connected">服务器已连接</div>
        ) : (
          <div className="server-status disconnected">
            服务器未连接 - 请确保后端服务器正在运行
          </div>
        )}
      </div>
      
      <main className="App-main">
        {isServerConnected !== false && (
          <>
            <div className="sidebar-container">
              <ChatList 
                onSelectChat={handleSelectChat}
                selectedChatId={selectedChatId}
                onChatDeleted={handleChatDeleted}
              />
            </div>
            <div className="chat-box-container">
              <ChatBox 
                username={username} 
                chatId={selectedChatId || undefined}
                onChatCreated={handleChatCreated}
              />
            </div>
          </>
        )}
      </main>
    </div>
  );
}

export default App;
