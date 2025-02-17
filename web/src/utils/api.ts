// src/utils/api.ts
export const backendUrl = 'http://localhost:8080'; // Adjust as needed

export async function registerUser(username: string, email: string, password: string) {
  const res = await fetch(`${backendUrl}/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password }),
  });
  if (!res.ok) {
    throw new Error(await res.text());
  }
  return await res.json();
}

export async function loginUser(email: string, password: string) {
  const res = await fetch(`${backendUrl}/v1/auth/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    throw new Error(await res.text());
  }
  // Option 2: The token is returned under the "data" property.
  return await res.json();
}