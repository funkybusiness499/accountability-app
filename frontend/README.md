# Accountability App Frontend

A modern web application built with Next.js, TypeScript, and Tailwind CSS that provides real-time accountability tracking features through WebSocket communication.

## Features

- Real-time chat functionality
- WebSocket-based communication
- Room-based messaging system
- Modern, responsive UI with Tailwind CSS
- Type-safe development with TypeScript

## Tech Stack

- **Framework**: Next.js 14.1.0
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context
- **Real-time Communication**: WebSocket

## Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js app directory
│   ├── components/       # React components
│   │   └── Chat.tsx     # Chat component
│   ├── contexts/        # React contexts
│   │   └── WebSocketContext.tsx
│   ├── services/        # Service classes
│   │   └── websocket.ts # WebSocket service
│   └── types/           # TypeScript types
│       └── websocket.ts # WebSocket types
├── public/              # Static files
└── tailwind.config.js   # Tailwind configuration
```

## Getting Started

1. **Install Dependencies**
   ```bash
   npm install
   ```

2. **Set Environment Variables**
   Create a `.env.local` file in the root directory:
   ```env
   NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8080/ws
   ```

3. **Run Development Server**
   ```bash
   npm run dev
   ```
   The application will be available at http://localhost:3000

4. **Build for Production**
   ```bash
   npm run build
   npm start
   ```

## WebSocket Integration

The application uses a custom WebSocket implementation with the following features:

- Automatic reconnection
- Room-based messaging
- Type-safe message handling
- Real-time presence updates

### Message Types

- `presence`: Room join/leave events
- `chat`: Text messages
- `task`: Task-related updates
- `notification`: System notifications

## Component Usage

### Chat Component

```tsx
import { Chat } from '@/components/Chat';

// In your page/component:
<Chat />
```

### WebSocket Context

```tsx
import { useWebSocket } from '@/contexts/WebSocketContext';

// In your component:
const { connected, currentRoom, messages, sendMessage } = useWebSocket();
```

## Available Scripts

- `npm run dev`: Start development server
- `npm run build`: Build production application
- `npm start`: Start production server
- `npm run lint`: Run ESLint

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

MIT License
