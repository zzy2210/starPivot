import React from 'react';
import UserProfile from '../User/UserProfile';
import './Header.css';

interface HeaderProps {
  onUsernameChange: (username: string) => void;
  username: string;
}

const Header: React.FC<HeaderProps> = ({ onUsernameChange, username }) => {
  return (
    <header className="header">
      <div className="header-content">
        <div className="header-left">
          <h1 className="app-title">StarPivot</h1>
        </div>
        <div className="header-right">
          <UserProfile onUsernameChange={onUsernameChange} username={username} />
        </div>
      </div>
    </header>
  );
};

export default Header; 