import React, { useState, useEffect } from 'react';
import './UserProfile.css';

interface UserProfileProps {
  onUsernameChange?: (username: string) => void;
  username?: string;
}

const UserProfile: React.FC<UserProfileProps> = ({ onUsernameChange, username: propUsername }) => {
  const [username, setUsername] = useState<string>(propUsername || '');
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [tempUsername, setTempUsername] = useState<string>('');

  // 当props中的username变化时更新状态
  useEffect(() => {
    if (propUsername && propUsername !== username) {
      setUsername(propUsername);
    }
  }, [propUsername, username]);

  // 组件加载时从localStorage获取用户名，仅在没有传入username props时使用
  useEffect(() => {
    if (!propUsername) {
      const savedUsername = localStorage.getItem('username');
      if (savedUsername) {
        setUsername(savedUsername);
        onUsernameChange && onUsernameChange(savedUsername);
      }
    }
  }, [onUsernameChange, propUsername]);

  // 保存用户名到localStorage
  const saveUsername = () => {
    if (tempUsername.trim()) {
      const newUsername = tempUsername.trim();
      localStorage.setItem('username', newUsername);
      setUsername(newUsername);
      setIsEditing(false);
      onUsernameChange && onUsernameChange(newUsername);
    }
  };

  // 开始编辑用户名
  const startEditing = () => {
    setTempUsername(username);
    setIsEditing(true);
  };

  // 取消编辑
  const cancelEditing = () => {
    setIsEditing(false);
  };

  // 处理输入变化
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setTempUsername(e.target.value);
  };

  // 处理表单提交
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    saveUsername();
  };

  return (
    <div className="user-profile">
      {!username && !isEditing ? (
        <div className="username-prompt">
          <p>请设置您的用户名：</p>
          <button onClick={startEditing} className="set-username-btn">
            设置用户名
          </button>
        </div>
      ) : !isEditing ? (
        <div className="username-display">
          <span className="username">{username}</span>
          <button onClick={startEditing} className="edit-username-btn">
            修改
          </button>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="username-form">
          <input
            type="text"
            value={tempUsername}
            onChange={handleInputChange}
            placeholder="输入用户名"
            autoFocus
            maxLength={20}
          />
          <div className="form-buttons">
            <button type="submit" disabled={!tempUsername.trim()}>
              保存
            </button>
            <button type="button" onClick={cancelEditing} className="cancel-btn">
              取消
            </button>
          </div>
        </form>
      )}
    </div>
  );
};

export default UserProfile; 