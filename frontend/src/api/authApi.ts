const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const TOKEN_KEY = "badminton_admin_token";

interface LoginResponse {
  token: string;
}

export async function login(
  email: string,
  password: string,
): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/auth/login`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    },
  );

  if (!response.ok) {
    throw new Error(await response.text());
  }

  const data: LoginResponse = await response.json();
  localStorage.setItem(TOKEN_KEY, data.token);
}

export function getAdminToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function logout(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export interface AdminProfile {
  id: string;
  name: string;
  email: string;
  role: string;
  eventId?: string;
}

export async function getAdminProfile(): Promise<AdminProfile> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/admin/me`,
    {
      headers: {
        Authorization: `Bearer ${getAdminToken()}`,
      },
    },
  );

  if (!response.ok) {
    throw new Error("Unable to load admin profile");
  }

  return response.json();
}