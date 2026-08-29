const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const SESSION_MARKER_KEY = "badminton_admin_session";

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
  sessionStorage.setItem(SESSION_MARKER_KEY, "active");
}

export function getAdminToken(): string | null {
  return null;
}

export function hasAdminSession(): boolean {
  return sessionStorage.getItem(SESSION_MARKER_KEY) === "active";
}

export function logout(): void {
	 sessionStorage.removeItem(SESSION_MARKER_KEY);
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