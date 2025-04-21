import { env } from '../config/env';

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  avatar?: string;
}

export class AuthService {
  private static readonly ACCESS_TOKEN_KEY = 'auth_access_token';
  private static readonly REFRESH_TOKEN_KEY = 'auth_refresh_token';
  private static readonly TOKEN_EXPIRY_KEY = 'auth_token_expiry';
  private static readonly USER_KEY = 'auth_user';
  private static readonly SESSION_KEY = 'auth_session';

  private baseUrl: string;
  private refreshPromise: Promise<AuthTokens> | null = null;
  private sessionTimeoutId?: NodeJS.Timeout;
  private readonly sessionTimeout = 30 * 60 * 1000; // 30 minutes

  constructor() {
    this.baseUrl = `${env.API_URL}/auth`;
    this.initializeSessionHandling();
  }

  private initializeSessionHandling(): void {
    // Restore session timeout if user is authenticated
    if (this.isAuthenticated()) {
      this.startSessionTimeout();
    }

    // Listen for storage events to handle multi-tab scenarios
    if (typeof window !== 'undefined') {
      window.addEventListener('storage', this.handleStorageChange);
    }
  }

  private handleStorageChange = (event: StorageEvent): void => {
    // Handle token changes from other tabs
    if (event.key === AuthService.ACCESS_TOKEN_KEY || event.key === AuthService.REFRESH_TOKEN_KEY) {
      if (!event.newValue) {
        // Token was removed in another tab
        this.clearSession();
      } else if (!event.oldValue && event.newValue) {
        // New token was added in another tab
        this.startSessionTimeout();
      }
    }
  };

  private startSessionTimeout(): void {
    this.clearSessionTimeout();
    this.sessionTimeoutId = setTimeout(() => {
      this.handleSessionTimeout();
    }, this.sessionTimeout);
    this.updateLastActivity();
  }

  private clearSessionTimeout(): void {
    if (this.sessionTimeoutId) {
      clearTimeout(this.sessionTimeoutId);
      this.sessionTimeoutId = undefined;
    }
  }

  private async handleSessionTimeout(): Promise<void> {
    const lastActivity = this.getLastActivity();
    const now = Date.now();

    if (lastActivity && now - lastActivity > this.sessionTimeout) {
      await this.logout();
      this.dispatchSessionEvent('session_expired');
    } else {
      this.startSessionTimeout();
    }
  }

  private updateLastActivity(): void {
    localStorage.setItem(AuthService.SESSION_KEY, Date.now().toString());
  }

  private getLastActivity(): number | null {
    const lastActivity = localStorage.getItem(AuthService.SESSION_KEY);
    return lastActivity ? parseInt(lastActivity, 10) : null;
  }

  private dispatchSessionEvent(type: 'session_expired' | 'session_recovered'): void {
    const event = new CustomEvent('auth_session_event', { detail: { type } });
    window.dispatchEvent(event);
  }

  async login(email: string, password: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const tokens: AuthTokens = await response.json();
    this.setTokens(tokens);
    this.startSessionTimeout();
  }

  async register(email: string, password: string, name: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, name }),
    });

    if (!response.ok) {
      throw new Error('Registration failed');
    }

    const tokens: AuthTokens = await response.json();
    this.setTokens(tokens);
    this.startSessionTimeout();
  }

  async logout(): Promise<void> {
    const refreshToken = this.getRefreshToken();
    if (refreshToken) {
      try {
        await fetch(`${this.baseUrl}/logout`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ refreshToken }),
        });
      } catch (error) {
        console.error('Logout request failed:', error);
      }
    }
    this.clearSession();
  }

  private clearSession(): void {
    this.clearTokens();
    this.clearSessionTimeout();
    localStorage.removeItem(AuthService.SESSION_KEY);
  }

  async getAccessToken(): Promise<string | null> {
    const token = localStorage.getItem(AuthService.ACCESS_TOKEN_KEY);
    const expiry = localStorage.getItem(AuthService.TOKEN_EXPIRY_KEY);

    if (!token || !expiry) {
      return null;
    }

    this.updateLastActivity();

    if (Date.now() >= parseInt(expiry, 10)) {
      return this.refreshAccessToken();
    }

    return token;
  }

  private async refreshAccessToken(): Promise<string | null> {
    if (this.refreshPromise) {
      return (await this.refreshPromise).accessToken;
    }

    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      this.clearSession();
      return null;
    }

    this.refreshPromise = this.performTokenRefresh(refreshToken);

    try {
      const tokens = await this.refreshPromise;
      this.setTokens(tokens);
      this.startSessionTimeout();
      return tokens.accessToken;
    } catch (error) {
      this.clearSession();
      return null;
    } finally {
      this.refreshPromise = null;
    }
  }

  private async performTokenRefresh(refreshToken: string): Promise<AuthTokens> {
    const response = await fetch(`${this.baseUrl}/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refreshToken }),
    });

    if (!response.ok) {
      throw new Error('Token refresh failed');
    }

    return response.json();
  }

  isAuthenticated(): boolean {
    return !!this.getRefreshToken();
  }

  getUser(): User | null {
    const userJson = localStorage.getItem(AuthService.USER_KEY);
    return userJson ? JSON.parse(userJson) : null;
  }

  private setTokens(tokens: AuthTokens): void {
    localStorage.setItem(AuthService.ACCESS_TOKEN_KEY, tokens.accessToken);
    localStorage.setItem(AuthService.REFRESH_TOKEN_KEY, tokens.refreshToken);
    localStorage.setItem(
      AuthService.TOKEN_EXPIRY_KEY,
      (Date.now() + tokens.expiresIn * 1000).toString()
    );
  }

  private getRefreshToken(): string | null {
    return localStorage.getItem(AuthService.REFRESH_TOKEN_KEY);
  }

  private clearTokens(): void {
    localStorage.removeItem(AuthService.ACCESS_TOKEN_KEY);
    localStorage.removeItem(AuthService.REFRESH_TOKEN_KEY);
    localStorage.removeItem(AuthService.TOKEN_EXPIRY_KEY);
    localStorage.removeItem(AuthService.USER_KEY);
  }

  async getAuthHeaders(): Promise<HeadersInit> {
    const token = await this.getAccessToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  onSessionEvent(callback: (event: { type: 'session_expired' | 'session_recovered' }) => void): () => void {
    const handler = (event: Event) => {
      const customEvent = event as CustomEvent<{ type: 'session_expired' | 'session_recovered' }>;
      callback(customEvent.detail);
    };

    window.addEventListener('auth_session_event', handler);
    return () => window.removeEventListener('auth_session_event', handler);
  }
} 