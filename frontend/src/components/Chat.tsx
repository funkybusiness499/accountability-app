import React, { useState, useEffect } from 'react';
import { useWebSocket } from '../contexts/WebSocketContext';
import { Message, MessageType, Room } from '../types/websocket';
import './Chat.css';

export const Chat: React.FC = () => {
  const [roomId, setRoomId] = useState<string>('');
  const [messageInput, setMessageInput] = useState<string>('');
  const { connected, currentRoom, messages, sendMessage, joinRoom } = useWebSocket();

  const handleJoinRoom = (e: React.FormEvent) => {
    e.preventDefault();
    if (roomId.trim()) {
      const room: Room = {
        id: roomId,
        name: `Room ${roomId}`,
        participants: []
      };
      joinRoom(room);
    }
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (messageInput.trim() && currentRoom) {
      sendMessage('chat', messageInput);
      setMessageInput('');
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h2>{currentRoom ? `Room: ${currentRoom.name}` : 'Join a Room'}</h2>
        <div>
          {connected ? (
            <span style={{ color: 'green' }}>Connected</span>
          ) : (
            <span style={{ color: 'red' }}>Disconnected</span>
          )}
        </div>
      </div>

      {!currentRoom ? (
        <form onSubmit={handleJoinRoom}>
          <div className="chat-input">
            <input
              type="text"
              value={roomId}
              onChange={(e) => setRoomId(e.target.value)}
              placeholder="Enter room ID"
              disabled={!connected}
            />
            <button type="submit" disabled={!connected || !roomId.trim()}>
              Join Room
            </button>
          </div>
        </form>
      ) : (
        <>
          <div className="chat-messages">
            {messages.filter(msg => msg.type === 'chat').map((msg, index) => (
              <div key={msg.id || index} className="message">
                <div>
                  <span className="sender">{msg.senderId}</span>
                  <span className="timestamp">
                    {new Date(msg.timestamp).toLocaleTimeString()}
                  </span>
                </div>
                <p className="content">{msg.data as string}</p>
              </div>
            ))}
          </div>

          <form onSubmit={handleSendMessage}>
            <div className="chat-input">
              <input
                type="text"
                value={messageInput}
                onChange={(e) => setMessageInput(e.target.value)}
                placeholder="Type a message..."
                disabled={!connected}
              />
              <button type="submit" disabled={!connected || !messageInput.trim()}>
                Send
              </button>
            </div>
          </form>
        </>
      )}
    </div>
  );
}; 