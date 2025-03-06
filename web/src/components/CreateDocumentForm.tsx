// src/components/CreateDocumentForm.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { backendUrl } from '../utils/api';

interface CreateDocumentFormProps {
  token: string;
}

const CreateDocumentForm: React.FC<CreateDocumentFormProps> = ({ token }) => {
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');
  const navigate = useNavigate();

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await fetch(`${backendUrl}/v1/document`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ content }), // only content is sent, docID generated on backend
      });
      if (!res.ok) {
        const errText = await res.text();
        setMessage(`Failed to create document: ${errText}`);
      } else {
        const data = await res.json();
        setMessage('Document created successfully!');
        // Redirect to the editor route with the generated docID.
        navigate(`/document/${data.doc_id}`);
      }
    } catch (error: any) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '500px', margin: 'auto', padding: '20px' }}>
      <h2>Create a New Document</h2>
      <form onSubmit={handleCreate}>
        {/* No document ID input needed */}
        <textarea
          placeholder="Initial Content (optional)"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          style={{ width: '100%', padding: '8px', marginBottom: '8px' }}
        />
        <button type="submit" disabled={loading} style={{ width: '100%', padding: '10px' }}>
          {loading ? 'Creating...' : 'Create Document'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
};

export default CreateDocumentForm;