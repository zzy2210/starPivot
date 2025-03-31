import React from 'react';
import ChatList from '../Chat/ChatList';
import './Sidebar.css';

interface SidebarProps {
  onSelectChat: (chatId: string) => void;
  selectedChatId: string | null;
}

const Sidebar: React.FC<SidebarProps> = ({ onSelectChat, selectedChatId }) => {
  return (
    <div className="sidebar">
      <ChatList onSelectChat={onSelectChat} selectedChatId={selectedChatId} />
    </div>
  );
};

export default Sidebar; 