// src/components/Editor.tsx
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Message } from '../types/message';

interface EditorProps {
  token: string;
  userID: string;
}

const Editor: React.FC<EditorProps> = ({ token, userID }) => {
  const { docID } = useParams<{ docID: string }>();
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<string>('Disconnected');
  const [content, setContent] = useState<string>('');

  // Construct the WebSocket URL including token and docID.
  const wsUrl = `ws://localhost:8080/v1/ws?docID=${encodeURIComponent(docID || '')}&token=${encodeURIComponent(token)}`;

  useEffect(() => {
    if (!docID) return;
    const socket = new WebSocket(wsUrl);
    setWs(socket);

    socket.onopen = () => {
      setConnectionStatus('Connected');
      console.log(`Connected to ${wsUrl}`);
    };

    socket.onmessage = (event) => {
      try {
        const msg: Message = JSON.parse(event.data);
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

    return () => {
      socket.close();
    };
  }, [wsUrl, docID]);

  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setContent(newContent);
    if (ws && ws.readyState === WebSocket.OPEN) {
      const message: Message = {
        type: 'update',
        docID: docID || '',
        position: 0, // Not used in this simple example.
        text: newContent,
        userID: userID,
        timestamp: new Date().toISOString(),
      };
      ws.send(JSON.stringify(message));
    }
  };

  return (
    <div style={{ padding: '20px' }}>
      <h1>DocCollab Editor</h1>
      <p>Status: {connectionStatus}</p>
      <p>Your User ID: {userID}</p>
      <p>Editing Document: {docID}</p>
      <textarea
        value={content}
        onChange={handleContentChange}
        style={{ width: '100%', height: '300px', padding: '10px', fontSize: '16px' }}
        placeholder="Start editing the shared document..."
      />
    </div>
  );
};

export default Editor;