import React, { useState, useEffect } from 'react';
import ChatBox from './components/Chat/ChatBox';
import { checkHealth } from './services/api';
import './App.css';

function App() {
  const [isServerConnected, setIsServerConnected] = useState<boolean | null>(null);

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
      <header className="App-header">
        {isServerConnected === null ? (
          <div className="server-status checking">正在连接服务器...</div>
        ) : isServerConnected ? (
          <div className="server-status connected">服务器已连接</div>
        ) : (
          <div className="server-status disconnected">
            服务器未连接 - 请确保后端服务器正在运行
          </div>
        )}
      </header>
      <main className="App-main">
        {isServerConnected !== false && <ChatBox />}
      </main>
    </div>
  );
}

export default App;
