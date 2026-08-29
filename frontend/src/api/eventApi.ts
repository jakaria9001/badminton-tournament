import type { EventInfo } from "../types/event";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL;

export async function getEvent(
  eventId: string,
): Promise<EventInfo> {

  const response = await fetch(
    `${API_BASE_URL}/api/v1/events/${eventId}`,
  );

  if (!response.ok) {
    throw new Error("Failed to load event");
  }

  return response.json();
}

export async function listEvents(): Promise<EventInfo[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/events`);

  if (!response.ok) {
    throw new Error("Failed to load events");
  }

  return response.json();
}

export async function listAdminEvents(): Promise<EventInfo[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/admin/superadmin/events`, {
	credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to load events");
  }

  return response.json();
}