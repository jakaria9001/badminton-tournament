const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

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
      credentials: "include",
      body: JSON.stringify({ email, password }),
    },
  );

  if (!response.ok) {
    throw new Error(await response.text());
  }
}

export function logout(): void {
  void fetch(`${API_BASE_URL}/api/v1/auth/logout`, {
    method: "POST",
    credentials: "include",
    keepalive: true,
  });
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
      credentials: "include",
    },
  );

  if (!response.ok) {
    throw new Error("Unable to load admin profile");
  }

  return response.json();
}