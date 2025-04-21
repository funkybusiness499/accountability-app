import { Room, User } from '../types/websocket';
import { env } from '../config/env';
import { AuthService } from './auth';

export class RoomService {
  private baseUrl: string;
  private authService: AuthService;

  constructor() {
    this.baseUrl = `${env.API_URL}/rooms`;
    this.authService = new AuthService();
  }

  async createRoom(name: string): Promise<Room> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...headers,
      },
      body: JSON.stringify({ name }),
    });

    if (!response.ok) {
      throw new Error('Failed to create room');
    }

    return response.json();
  }

  async getRooms(): Promise<Room[]> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(this.baseUrl, {
      headers,
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch rooms');
    }

    return response.json();
  }

  async getRoom(roomId: string): Promise<Room> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(`${this.baseUrl}/${roomId}`, {
      headers,
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch room');
    }

    return response.json();
  }

  async getRoomParticipants(roomId: string): Promise<User[]> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(`${this.baseUrl}/${roomId}/participants`, {
      headers,
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch room participants');
    }

    return response.json();
  }

  async joinRoom(roomId: string): Promise<void> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(`${this.baseUrl}/${roomId}/join`, {
      method: 'POST',
      headers,
    });

    if (!response.ok) {
      throw new Error('Failed to join room');
    }
  }

  async leaveRoom(roomId: string): Promise<void> {
    const headers = await this.authService.getAuthHeaders();
    const response = await fetch(`${this.baseUrl}/${roomId}/leave`, {
      method: 'POST',
      headers,
    });

    if (!response.ok) {
      throw new Error('Failed to leave room');
    }
  }
} 