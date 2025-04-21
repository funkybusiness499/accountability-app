import React, { createContext, useContext, useEffect, useState } from 'react';
import { WebSocketService } from '../services/websocket';
import { RoomService } from '../services/room';
import { Message, Room, WebSocketContextType, WebSocketState } from '../types/websocket';

const WebSocketContext = createContext<WebSocketContextType | null>(null);

const initialState: WebSocketState = {
  connected: false,
  currentRoom: null,
  messages: [],
  error: null,
};

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, setState] = useState<WebSocketState>(initialState);
  const [wsService] = useState(() => new WebSocketService());
  const [roomService] = useState(() => new RoomService());

  useEffect(() => {
    const handleMessage = (message: Message) => {
      setState(prev => ({
        ...prev,
        messages: [...prev.messages, message],
      }));
    };

    const handleStatus = (connected: boolean) => {
      setState(prev => ({ ...prev, connected, error: null }));
    };

    const handleError = (error: string) => {
      setState(prev => ({ ...prev, error }));
    };

    wsService.onMessage(handleMessage);
    wsService.onStatusChange(handleStatus);
    wsService.onError(handleError);

    return () => {
      wsService.disconnect();
    };
  }, [wsService]);

  const connect = () => {
    wsService.connect();
  };

  const disconnect = () => {
    wsService.disconnect();
    setState(initialState);
  };

  const joinRoom = async (room: Room) => {
    try {
      await roomService.joinRoom(room.id);
      wsService.sendMessage('presence', { action: 'join', roomId: room.id });
      setState(prev => ({ ...prev, currentRoom: room }));
    } catch (error) {
      setState(prev => ({ ...prev, error: 'Failed to join room' }));
    }
  };

  const leaveRoom = async () => {
    if (!state.currentRoom) return;

    try {
      await roomService.leaveRoom(state.currentRoom.id);
      wsService.sendMessage('presence', { action: 'leave', roomId: state.currentRoom.id });
      setState(prev => ({ ...prev, currentRoom: null }));
    } catch (error) {
      setState(prev => ({ ...prev, error: 'Failed to leave room' }));
    }
  };

  const sendMessage = (type: Message['type'], data: any) => {
    if (!state.currentRoom) {
      setState(prev => ({ ...prev, error: 'Not connected to any room' }));
      return;
    }

    wsService.sendMessage(type, data, state.currentRoom.id);
  };

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