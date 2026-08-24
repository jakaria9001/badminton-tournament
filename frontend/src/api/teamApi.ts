export interface TeamPlayer {
  id: string;
  name: string;
}

export interface Team {
  id: string;
  teamName: string;
  player1: TeamPlayer;
  player2: TeamPlayer;
  status: string;
  createdAt: string;
}

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL;

export async function getTeams(
  eventId: string,
): Promise<Team[]> {

  const response = await fetch(
    `${API_BASE_URL}/api/v1/events/${eventId}/teams`,
  );

  if (!response.ok) {
    throw new Error("Failed to load teams");
  }

  const data = await response.json();

  return data.teams ?? [];
}