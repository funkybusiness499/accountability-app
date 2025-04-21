import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { Message, MessageType, Room, WebSocketContextType, WebSocketState } from '../types/websocket';

const WebSocketContext = createContext<WebSocketContextType | null>(null);

const WEBSOCKET_URL = process.env.REACT_APP_WEBSOCKET_URL || 'ws://localhost:8080/ws';

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const ws = useRef<WebSocket | null>(null);
  const [state, setState] = useState<WebSocketState>({
    connected: false,
    currentRoom: null,
    messages: [],
    error: null,
  });

  const connect = () => {
    try {
      ws.current = new WebSocket(WEBSOCKET_URL);
      
      ws.current.onopen = () => {
        setState(prev => ({ ...prev, connected: true, error: null }));
      };

      ws.current.onclose = () => {
        setState(prev => ({ ...prev, connected: false }));
      };

      ws.current.onerror = (error) => {
        setState(prev => ({ ...prev, error: 'WebSocket connection error' }));
      };

      ws.current.onmessage = (event) => {
        try {
          const message: Message = JSON.parse(event.data);
          setState(prev => ({
            ...prev,
            messages: [...prev.messages, message],
          }));
        } catch (error) {
          setState(prev => ({ ...prev, error: 'Invalid message format' }));
        }
      };
    } catch (error) {
      setState(prev => ({ ...prev, error: 'Failed to connect to WebSocket' }));
    }
  };

  const disconnect = () => {
    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
  };

  const joinRoom = (room: Room) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = {
        type: 'presence' as MessageType,
        data: { action: 'join', roomId: room.id }
      };
      ws.current.send(JSON.stringify(message));
      setState(prev => ({ ...prev, currentRoom: room }));
    }
  };

  const leaveRoom = () => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN && state.currentRoom) {
      const message = {
        type: 'presence' as MessageType,
        data: { action: 'leave', roomId: state.currentRoom.id }
      };
      ws.current.send(JSON.stringify(message));
      setState(prev => ({ ...prev, currentRoom: null }));
    }
  };

  const sendMessage = (type: MessageType, data: any) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN && state.currentRoom) {
      const message = {
        type,
        roomId: state.currentRoom.id,
        data
      };
      ws.current.send(JSON.stringify(message));
    }
  };

  useEffect(() => {
    return () => {
      disconnect();
    };
  }, []);

  const value: WebSocketContextType = {
    ...state,
    connect,
    disconnect,
    joinRoom,
    leaveRoom,
    sendMessage,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
}; 