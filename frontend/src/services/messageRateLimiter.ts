interface MessageBucket {
  tokens: number;
  lastRefill: number;
}

export class MessageRateLimiter {
  private static readonly STORAGE_KEY = 'message_rate_limit';
  private bucket: MessageBucket;

  constructor(
    private maxTokens: number = 60, // Maximum messages per minute
    private refillRate: number = 1, // Tokens to add per second
    private refillInterval: number = 1000 // Refill interval in milliseconds
  ) {
    this.bucket = this.loadBucket();
    this.startRefillTimer();
  }

  canSendMessage(): boolean {
    this.refillTokens();
    return this.bucket.tokens > 0;
  }

  consumeToken(): void {
    if (this.bucket.tokens > 0) {
      this.bucket.tokens--;
      this.saveBucket();
    }
  }

  getTimeUntilNextMessage(): number {
    if (this.bucket.tokens > 0) return 0;
    
    const timeSinceLastRefill = Date.now() - this.bucket.lastRefill;
    const tokensToAdd = Math.floor(timeSinceLastRefill / this.refillInterval) * this.refillRate;
    
    if (tokensToAdd > 0) {
      this.refillTokens();
      return this.bucket.tokens > 0 ? 0 : this.refillInterval;
    }
    
    return this.refillInterval - (timeSinceLastRefill % this.refillInterval);
  }

  private refillTokens(): void {
    const now = Date.now();
    const timeSinceLastRefill = now - this.bucket.lastRefill;
    const tokensToAdd = Math.floor(timeSinceLastRefill / this.refillInterval) * this.refillRate;

    if (tokensToAdd > 0) {
      this.bucket.tokens = Math.min(this.maxTokens, this.bucket.tokens + tokensToAdd);
      this.bucket.lastRefill = now;
      this.saveBucket();
    }
  }

  private startRefillTimer(): void {
    if (typeof window !== 'undefined') {
      setInterval(() => {
        this.refillTokens();
      }, this.refillInterval);
    }
  }

  private loadBucket(): MessageBucket {
    const stored = localStorage.getItem(MessageRateLimiter.STORAGE_KEY);
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch {
        return this.getInitialBucket();
      }
    }
    return this.getInitialBucket();
  }

  private getInitialBucket(): MessageBucket {
    return {
      tokens: this.maxTokens,
      lastRefill: Date.now(),
    };
  }

  private saveBucket(): void {
    localStorage.setItem(MessageRateLimiter.STORAGE_KEY, JSON.stringify(this.bucket));
  }
} 