import { env } from '../config/env';
import { Message, MessageType } from '../types/websocket';

export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 5;
  private reconnectTimeout: number = 1000;
  private heartbeatInterval: number = 30000;
  private heartbeatTimer?: NodeJS.Timeout;

  private messageHandlers: ((message: Message) => void)[] = [];
  private statusHandlers: ((connected: boolean) => void)[] = [];
  private errorHandlers: ((error: string) => void)[] = [];

  constructor(private authToken?: string) {}

  connect(): void {
    try {
      const url = new URL(env.WEBSOCKET_URL);
      if (this.authToken) {
        url.searchParams.append('token', this.authToken);
      }

      this.ws = new WebSocket(url.toString());
      this.setupEventListeners();
    } catch (error) {
      this.handleError('Failed to establish WebSocket connection');
    }
  }

  private setupEventListeners(): void {
    if (!this.ws) return;

    this.ws.onopen = () => {
      console.log('WebSocket connection established');
      this.reconnectAttempts = 0;
      this.startHeartbeat();
      this.notifyStatusChange(true);
    };

    this.ws.onclose = () => {
      console.log('WebSocket connection closed');
      this.stopHeartbeat();
      this.notifyStatusChange(false);
      this.attemptReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.handleError('WebSocket connection error');
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as Message;
        if (message.type === 'ping') {
          this.sendPong();
        } else {
          this.notifyMessageReceived(message);
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
        this.handleError('Invalid message format');
      }
    };
  }

  private startHeartbeat(): void {
    this.heartbeatTimer = setInterval(() => {
      this.sendMessage('ping', null);
    }, this.heartbeatInterval);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = undefined;
    }
  }

  private sendPong(): void {
    this.sendMessage('pong', null);
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.handleError('Maximum reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

    setTimeout(() => {
      this.connect();
    }, this.reconnectTimeout * this.reconnectAttempts);
  }

  disconnect(): void {
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  sendMessage(type: MessageType, data: any, roomId?: string): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.handleError('WebSocket is not connected');
      return;
    }

    const message = {
      type,
      data,
      roomId,
      timestamp: Date.now(),
    };

    try {
      this.ws.send(JSON.stringify(message));
    } catch (error) {
      console.error('Failed to send message:', error);
      this.handleError('Failed to send message');
    }
  }

  onMessage(handler: (message: Message) => void): () => void {
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter(h => h !== handler);
    };
  }

  onStatusChange(handler: (connected: boolean) => void): () => void {
    this.statusHandlers.push(handler);
    return () => {
      this.statusHandlers = this.statusHandlers.filter(h => h !== handler);
    };
  }

  onError(handler: (error: string) => void): () => void {
    this.errorHandlers.push(handler);
    return () => {
      this.errorHandlers = this.errorHandlers.filter(h => h !== handler);
    };
  }

  private notifyMessageReceived(message: Message): void {
    if (env.DEBUG_MODE) {
      console.log('Received message:', message);
    }
    this.messageHandlers.forEach(handler => handler(message));
  }

  private notifyStatusChange(connected: boolean): void {
    this.statusHandlers.forEach(handler => handler(connected));
  }

  private handleError(error: string): void {
    if (env.DEBUG_MODE) {
      console.error('WebSocket error:', error);
    }
    this.errorHandlers.forEach(handler => handler(error));
  }
} 