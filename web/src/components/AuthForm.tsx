// src/components/AuthForm.tsx
import React, { useState } from 'react';
import { loginUser, registerUser } from '../utils/api';

export type AuthMode = 'login' | 'register';

interface AuthFormProps {
  mode: AuthMode;
  onAuthSuccess: (token: string) => void;
  switchMode: (mode: AuthMode) => void;
}

const AuthForm: React.FC<AuthFormProps> = ({ mode, onAuthSuccess, switchMode }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState(''); // only used for registration
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      if (mode === 'register') {
        await registerUser(username, email, password);
        alert('Registration successful! Please log in.');
        switchMode('login');
      } else {
        const data = await loginUser(email, password);
        // Extract the token from the "data" field.
        onAuthSuccess(data.data);
      }
    } catch (err: any) {
      alert(`${mode === 'register' ? 'Registration' : 'Login'} failed: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: 'auto', padding: '20px' }}>
      <h2>{mode === 'login' ? 'Login' : 'Register'}</h2>
      <form onSubmit={handleSubmit}>
        {mode === 'register' && (
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={{ width: '100%', padding: '8px', marginBottom: '8px' }}
            required
          />
        )}
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          style={{ width: '100%', padding: '8px', marginBottom: '8px' }}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={{ width: '100%', padding: '8px', marginBottom: '8px' }}
          required
        />
        <button type="submit" style={{ width: '100%', padding: '10px' }} disabled={loading}>
          {loading ? 'Loading...' : mode === 'login' ? 'Login' : 'Register'}
        </button>
      </form>
      <div style={{ marginTop: '10px' }}>
        {mode === 'login' ? (
          <p>
            Don't have an account?{' '}
            <button onClick={() => switchMode('register')}>Register</button>
          </p>
        ) : (
          <p>
            Already have an account?{' '}
            <button onClick={() => switchMode('login')}>Login</button>
          </p>
        )}
      </div>
    </div>
  );
};

export default AuthForm;