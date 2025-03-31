import React from 'react';
import UserProfile from '../User/UserProfile';
import './Header.css';

interface HeaderProps {
  onUsernameChange: (username: string) => void;
}

const Header: React.FC<HeaderProps> = ({ onUsernameChange }) => {
  return (
    <header className="header">
      <div className="header-content">
        <div className="header-left">
          <h1 className="app-title">StarPivot</h1>
        </div>
        <div className="header-right">
          <UserProfile onUsernameChange={onUsernameChange} />
        </div>
      </div>
    </header>
  );
};

export default Header; 