
const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL;

export interface PlayerInput {
  name: string;
  phone: string;
}

export interface RegistrationRequest {
  player1: PlayerInput;
  player2: PlayerInput;
  teamName?: string;
}

export interface RegistrationResponse {
  registrationId: string;
  status: string;
}

export async function registerTeam(
  eventId: string,
  request: RegistrationRequest,
): Promise<RegistrationResponse> {

  const response = await fetch(
    `${API_BASE_URL}/api/v1/events/${eventId}/registrations`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    },
  );

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Registration failed");
  }

  return response.json();
}