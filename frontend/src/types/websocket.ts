export type MessageType = 'presence' | 'chat' | 'task' | 'notification' | 'ping' | 'pong';

export interface User {
  id: string;
  name: string;
  avatar?: string;
}

export interface Message {
  id?: string;
  type: MessageType;
  roomId?: string;
  senderId?: string;
  data: any;
  timestamp: number;
}

export interface Room {
  id: string;
  name: string;
  participants: User[];
}

export interface WebSocketState {
  connected: boolean;
  currentRoom: Room | null;
  messages: Message[];
  error: string | null;
}

export interface WebSocketContextType extends WebSocketState {
  connect: () => void;
  disconnect: () => void;
  joinRoom: (room: Room) => void;
  leaveRoom: () => void;
  sendMessage: (type: MessageType, data: any) => void;
} 