interface EnvConfig {
  API_URL: string;
  WEBSOCKET_URL: string;
  ENABLE_ANALYTICS: boolean;
  ENABLE_AUTH: boolean;
  DEBUG_MODE: boolean;
  API_TIMEOUT: number;
}

const requiredEnvVars = [
  'NEXT_PUBLIC_API_URL',
  'NEXT_PUBLIC_WEBSOCKET_URL',
  'NEXT_PUBLIC_ENABLE_ANALYTICS',
  'NEXT_PUBLIC_ENABLE_AUTH',
  'NEXT_PUBLIC_DEBUG_MODE',
  'NEXT_PUBLIC_API_TIMEOUT',
] as const;

function validateEnv(): EnvConfig {
  for (const envVar of requiredEnvVars) {
    if (!(envVar in process.env)) {
      throw new Error(`Missing required environment variable: ${envVar}`);
    }
  }

  return {
    API_URL: process.env.NEXT_PUBLIC_API_URL!,
    WEBSOCKET_URL: process.env.NEXT_PUBLIC_WEBSOCKET_URL!,
    ENABLE_ANALYTICS: process.env.NEXT_PUBLIC_ENABLE_ANALYTICS === 'true',
    ENABLE_AUTH: process.env.NEXT_PUBLIC_ENABLE_AUTH === 'true',
    DEBUG_MODE: process.env.NEXT_PUBLIC_DEBUG_MODE === 'true',
    API_TIMEOUT: parseInt(process.env.NEXT_PUBLIC_API_TIMEOUT!, 10),
  };
}

export const env = validateEnv(); 