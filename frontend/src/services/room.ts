import { Room, User } from '../types/websocket';
import { env } from '../config/env';

export class RoomService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = `${env.API_URL}/rooms`;
  }

  async createRoom(name: string): Promise<Room> {
    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name }),
    });

    if (!response.ok) {
      throw new Error('Failed to create room');
    }

    return response.json();
  }

  async getRooms(): Promise<Room[]> {
    const response = await fetch(this.baseUrl);
    
    if (!response.ok) {
      throw new Error('Failed to fetch rooms');
    }

    return response.json();
  }

  async getRoom(roomId: string): Promise<Room> {
    const response = await fetch(`${this.baseUrl}/${roomId}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch room');
    }

    return response.json();
  }

  async getRoomParticipants(roomId: string): Promise<User[]> {
    const response = await fetch(`${this.baseUrl}/${roomId}/participants`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch room participants');
    }

    return response.json();
  }

  async joinRoom(roomId: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${roomId}/join`, {
      method: 'POST',
    });

    if (!response.ok) {
      throw new Error('Failed to join room');
    }
  }

  async leaveRoom(roomId: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${roomId}/leave`, {
      method: 'POST',
    });

    if (!response.ok) {
      throw new Error('Failed to leave room');
    }
  }
} 