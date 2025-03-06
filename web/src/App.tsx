// src/App.tsx
import React, { useState } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import AuthForm, { AuthMode } from './components/AuthForm';
import CreateDocumentForm from './components/CreateDocumentForm';
import Editor from './components/Editor';

const App: React.FC = () => {
  // Holds the JWT token. If token exists, user is authenticated.
  const [token, setToken] = useState<string>('');
  // For simplicity, generate a random user ID on initial login.
  // In a real app, you might decode this from the token.
  const [userID] = useState<string>(() => Math.random().toString(36).substring(2, 10));
  // Manage authentication mode: either 'login' or 'register'
  const [authMode, setAuthMode] = useState<AuthMode>('login');

  return (
    <Routes>
      {/* If not authenticated, route to /auth */}
      {!token ? (
        <Route
          path="/*"
          element={
            <AuthForm
              mode={authMode}
              onAuthSuccess={(t: string) => setToken(t)}
              switchMode={(mode: AuthMode) => setAuthMode(mode)}
            />
          }
        />
      ) : (
        <>
          {/* Default route when logged in goes to create document */}
          <Route path="/" element={<Navigate to="/document/create" />} />
          {/* Route for creating a new document */}
          <Route
            path="/document/create"
            element={<CreateDocumentForm token={token} />}
          />
          {/* Editor route, expecting a generated docID */}
          <Route
            path="/document/:docID"
            element={<Editor token={token} userID={userID} />}
          />
        </>
      )}
    </Routes>
  );
};

export default App;