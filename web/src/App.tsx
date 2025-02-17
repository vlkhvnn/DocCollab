// src/App.tsx
import React, { useState } from 'react';
import AuthForm, { AuthMode } from './components/AuthForm';
import Editor from './components/Editor';

const App: React.FC = () => {
  // Holds the JWT token. If token exists, user is authenticated.
  const [token, setToken] = useState<string>('');
  // For simplicity, generate a random user ID on initial login.
  // In a real app, you might decode this from the token.
  const [userID] = useState<string>(() => Math.random().toString(36).substring(2, 10));
  // Manage authentication mode: 'login' or 'register'
  const [authMode, setAuthMode] = useState<AuthMode>('login');

  // If no token, show the auth form.
  if (!token) {
    return (
      <div className="App">
        <AuthForm
          mode={authMode}
          onAuthSuccess={(t: string) => setToken(t)}
          switchMode={(mode: AuthMode) => setAuthMode(mode)}
        />
      </div>
    );
  }

  // Once authenticated, show the collaborative editor.
  return <Editor token={token} userID={userID} />;
};

export default App;