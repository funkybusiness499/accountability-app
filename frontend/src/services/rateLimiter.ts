interface RateLimitConfig {
  maxAttempts: number;
  timeWindow: number;
  initialDelay: number;
  maxDelay: number;
  jitter: boolean;
}

interface RateLimitState {
  attempts: number;
  resetTime: number;
  nextAttemptTime: number;
  currentDelay: number;
}

export class RateLimiter {
  private static readonly STATE_PREFIX = 'rate_limit_';
  private config: RateLimitConfig;
  private state: RateLimitState;

  constructor(
    private key: string,
    config?: Partial<RateLimitConfig>
  ) {
    this.config = {
      maxAttempts: config?.maxAttempts ?? 5,
      timeWindow: config?.timeWindow ?? 60000, // 1 minute
      initialDelay: config?.initialDelay ?? 1000, // 1 second
      maxDelay: config?.maxDelay ?? 32000, // 32 seconds
      jitter: config?.jitter ?? true,
    };

    this.state = this.loadState();
  }

  canAttempt(): boolean {
    const now = Date.now();

    // Reset state if time window has passed
    if (now >= this.state.resetTime) {
      this.resetState();
      return true;
    }

    // Check if we're still in backoff period
    if (now < this.state.nextAttemptTime) {
      return false;
    }

    // Check if we've exceeded max attempts
    return this.state.attempts < this.config.maxAttempts;
  }

  recordAttempt(success: boolean): void {
    const now = Date.now();

    if (success) {
      // On success, gradually reduce the delay
      this.state.currentDelay = Math.max(
        this.config.initialDelay,
        this.state.currentDelay / 2
      );
    } else {
      // On failure, increment attempts and apply exponential backoff
      this.state.attempts++;
      this.state.currentDelay = Math.min(
        this.config.maxDelay,
        this.state.currentDelay * 2 || this.config.initialDelay
      );

      // Add jitter if enabled
      if (this.config.jitter) {
        this.state.currentDelay *= 0.5 + Math.random();
      }
    }

    this.state.nextAttemptTime = now + this.state.currentDelay;
    this.saveState();
  }

  getNextAttemptDelay(): number {
    const now = Date.now();
    return Math.max(0, this.state.nextAttemptTime - now);
  }

  isRateLimited(): boolean {
    return this.state.attempts >= this.config.maxAttempts;
  }

  private resetState(): void {
    this.state = {
      attempts: 0,
      resetTime: Date.now() + this.config.timeWindow,
      nextAttemptTime: 0,
      currentDelay: this.config.initialDelay,
    };
    this.saveState();
  }

  private loadState(): RateLimitState {
    const stored = localStorage.getItem(this.getStorageKey());
    if (stored) {
      try {
        const state = JSON.parse(stored) as RateLimitState;
        // If stored state has expired, return fresh state
        if (Date.now() >= state.resetTime) {
          return this.getInitialState();
        }
        return state;
      } catch {
        return this.getInitialState();
      }
    }
    return this.getInitialState();
  }

  private getInitialState(): RateLimitState {
    return {
      attempts: 0,
      resetTime: Date.now() + this.config.timeWindow,
      nextAttemptTime: 0,
      currentDelay: this.config.initialDelay,
    };
  }

  private saveState(): void {
    localStorage.setItem(this.getStorageKey(), JSON.stringify(this.state));
  }

  private getStorageKey(): string {
    return `${RateLimiter.STATE_PREFIX}${this.key}`;
  }
} 