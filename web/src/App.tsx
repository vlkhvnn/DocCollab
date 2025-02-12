import React, { useEffect, useState } from 'react';
import './App.css';

// Define the message structure that matches the backend.
interface Message {
  type: string;      // "update" for client updates, "sync" for server sync messages
  docID: string;
  position: number;  // not used in this simple example
  text: string;
  userID: string;
  timestamp: string;
}

// Use a fixed document ID for demonstration.
const docID = 'doc123';
// Change the URL if your backend is hosted elsewhere.
const wsUrl = `ws://localhost:8080/v1/ws?docID=${encodeURIComponent(docID)}`;

const App: React.FC = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<string>('Disconnected');
  const [content, setContent] = useState<string>(''); // shared document content
  // Generate a random user ID when the app first loads.
  const [userId] = useState<string>(() => Math.random().toString(36).substring(2, 10));

  useEffect(() => {
    // Create the WebSocket connection.
    const socket = new WebSocket(wsUrl);
    setWs(socket);

    socket.onopen = () => {
      setConnectionStatus('Connected');
      console.log(`Connected to ${wsUrl}`);
    };

    socket.onmessage = (event) => {
      try {
        const msg: Message = JSON.parse(event.data);
        // When a sync message is received, update the shared content.
        if (msg.type === 'sync') {
          console.log('Sync received:', msg);
          setContent(msg.text);
        }
      } catch (err) {
        console.error('Error parsing message:', err);
      }
    };

    socket.onerror = (err) => {
      console.error('WebSocket error:', err);
    };

    socket.onclose = () => {
      setConnectionStatus('Disconnected');
      console.log('WebSocket disconnected');
    };

    // Clean up on component unmount.
    return () => {
      socket.close();
    };
  }, []);

  // When the user types, update the local state and send an update to the backend.
  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setContent(newContent);

    // Send an "update" message with the full content.
    if (ws && ws.readyState === WebSocket.OPEN) {
      const message: Message = {
        type: 'update',
        docID: docID,
        position: 0, // not used here
        text: newContent,
        userID: userId,
        timestamp: new Date().toISOString(),
      };
      ws.send(JSON.stringify(message));
    }
  };

  return (
    <div className="App" style={{ padding: '20px' }}>
      <h1>DocCollab Editor</h1>
      <p>Status: {connectionStatus}</p>
      <p>Your User ID: {userId}</p>
      <textarea
        value={content}
        onChange={handleContentChange}
        style={{ width: '100%', height: '300px', padding: '10px', fontSize: '16px' }}
        placeholder="Start editing the shared document..."
      />
    </div>
  );
};

export default App;