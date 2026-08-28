import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getAdminToken, logout } from "../../api/authApi";
import PublicFooter from "../../components/PublicFooter";
import PublicHeader from "../../components/PublicHeader";

const EVENT_ID = "00000000-0000-0000-0000-000000000002";
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

interface Round {
  id: string;
  roundNumber: number;
  roundName: string;
  pairingMethod: string;
  status: string;
}

interface Team {
  id: string;
  teamName: string;
  player1: string;
  player2: string;
}

interface Match {
  id: string;
  roundId: string;
  team1: Team | null;
  team2: Team | null;
  status: string;
}

function Draw() {
  const navigate = useNavigate();
  const [rounds, setRounds] = useState<Round[]>([]);
  const [teams, setTeams] = useState<Record<string, Team[]>>({});
  const [matches, setMatches] = useState<Record<string, Match[]>>({});
  const [selectedTeams, setSelectedTeams] = useState<Record<string, string>>({});
  const [newRoundPairingMethod, setNewRoundPairingMethod] = useState("RANDOM");
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState<string | null>(null);
  const [message, setMessage] = useState("");

  const request = useCallback(async (url: string, options: RequestInit = {}) => {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getAdminToken()}`,
        ...options.headers,
      },
    });

    if (response.status === 401 || response.status === 403) {
      logout();
      navigate("/admin/login", { replace: true });
      throw new Error("Session expired");
    }

    if (!response.ok) {
      throw new Error(await response.text());
    }

    return response;
  }, [navigate]);

  const loadRounds = useCallback(async () => {
    const response = await request(`/api/v1/admin/events/${EVENT_ID}/rounds`);
    const data = await response.json();
    setRounds(data ?? []);

    const matchesResponse = await request(`/api/v1/events/${EVENT_ID}/matches`);
    const matchData: Match[] = (await matchesResponse.json()) ?? [];
    const matchesByRound = matchData.reduce<Record<string, Match[]>>((grouped, match) => {
      grouped[match.roundId] = [...(grouped[match.roundId] ?? []), match];
      return grouped;
    }, {});
    setMatches(matchesByRound);

    const available = await Promise.all((data ?? []).map(async (round: Round) => {
      const teamResponse = await request(
        `/api/v1/admin/events/${EVENT_ID}/rounds/${round.id}/available-teams`,
      );
      return [round.id, (await teamResponse.json()) ?? []] as const;
    }));
    setTeams(Object.fromEntries(available));
  }, [request]);

  useEffect(() => {
    void loadRounds().catch((error: unknown) => {
      if (error instanceof Error && error.message !== "Session expired") {
        setMessage(error.message);
      }
    }).finally(() => setLoading(false));
  }, [loadRounds]);

  async function createRound() {
    setBusy("create");
    setMessage("");
    try {
      await request(`/api/v1/admin/events/${EVENT_ID}/rounds`, {
        method: "POST",
        body: JSON.stringify({
          roundNumber: rounds.length + 1,
          roundName: `ROUND_${rounds.length + 1}`,
          pairingMethod: newRoundPairingMethod,
        }),
      });
      await loadRounds();
      setMessage("Round created.");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Unable to create round");
    } finally {
      setBusy(null);
    }
  }

  async function generate(round: Round) {
    setBusy(round.id);
    setMessage("");
    try {
      await request(`/api/v1/admin/rounds/${round.id}/generate`, { method: "POST" });
      await loadRounds();
      setMessage("Draw generated.");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Unable to generate draw");
    } finally {
      setBusy(null);
    }
  }

  async function createMatch(round: Round) {
    const team1Id = selectedTeams[`${round.id}:team1`];
    const team2Id = selectedTeams[`${round.id}:team2`];
    if (!team1Id || !team2Id) {
      setMessage("Select two teams first.");
      return;
    }

    setBusy(round.id);
    try {
      await request(`/api/v1/admin/rounds/${round.id}/matches`, {
        method: "POST",
        body: JSON.stringify({ team1Id, team2Id }),
      });
      await loadRounds();
      setMessage("Match created.");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Unable to create match");
    } finally {
      setBusy(null);
    }
  }

  async function lock(round: Round) {
    setBusy(round.id);
    try {
      await request(`/api/v1/admin/rounds/${round.id}/lock`, { method: "POST" });
      await loadRounds();
      setMessage("Round locked.");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Unable to lock round");
    } finally {
      setBusy(null);
    }
  }

  if (loading) {
    return (
      <div className="flex min-h-screen flex-col bg-slate-950 text-white">
        <PublicHeader />
        <main className="flex flex-1 items-center justify-center px-5 py-10">
          Loading draw...
        </main>
        <PublicFooter />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900">
      <PublicHeader />
      <main className="flex-1 px-4 py-8">
        <div className="mx-auto max-w-5xl">
          <div className="mb-6 flex flex-wrap items-end justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.18em] text-slate-500">Admin</p>
              <h1 className="mt-2 text-3xl font-bold text-slate-950">Tournament Draw</h1>
            </div>
            <div className="flex gap-2">
              <button className="rounded-lg border border-slate-300 bg-white px-4 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-50" onClick={() => navigate("/admin/registrations")} type="button">Registrations</button>
              <button className="rounded-lg border border-red-200 bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 transition hover:bg-red-100" onClick={() => { logout(); navigate("/admin/login", { replace: true }); }} type="button">Sign out</button>
            </div>
          </div>

          {message && <p className="mb-4 rounded-lg bg-white p-4 text-sm font-medium text-slate-700 shadow-sm">{message}</p>}

          <div className="mb-6 flex flex-wrap items-end justify-end gap-3">
            <label className="flex flex-col gap-1 text-sm font-semibold text-slate-700">
              Pairing method
              <select
                className="rounded-lg border border-slate-300 bg-white px-3 py-2 font-normal"
                value={newRoundPairingMethod}
                onChange={(event) => setNewRoundPairingMethod(event.target.value)}
              >
                <option value="RANDOM">Random</option>
                <option value="MANUAL">Manual</option>
              </select>
            </label>
            <button className="rounded-lg bg-emerald-600 px-4 py-2 font-semibold text-white disabled:opacity-50" disabled={busy !== null} onClick={() => void createRound()} type="button">Create round</button>
          </div>

          <div className="space-y-4">
            {rounds.length === 0 && <div className="rounded-2xl bg-white p-8 text-center text-slate-500 shadow-sm">No rounds created yet.</div>}
            {rounds.map((round) => {
                const roundMatches = matches[round.id] ?? [];
                const usedTeamIds = new Set(
                  roundMatches.flatMap((match) => [match.team1?.id, match.team2?.id]).filter(Boolean),
                );
                const availableTeams = (teams[round.id] ?? []).filter((team, index, allTeams) => (
                  !usedTeamIds.has(team.id) && allTeams.findIndex((candidate) => candidate.id === team.id) === index
                ));
              return (
                <section className="rounded-2xl bg-white p-6 shadow-sm" key={round.id}>
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <h2 className="text-xl font-bold text-slate-900">{round.roundName}</h2>
                    <p className="text-sm text-slate-500">Round {round.roundNumber} · {round.pairingMethod}</p>
                  </div>
                  <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-bold text-slate-700">{round.status}</span>
                </div>

                {round.status === "OPEN" && (
                  <div className="mt-5 flex flex-wrap gap-3">
                    {round.pairingMethod === "RANDOM" && (
                      <button className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50" disabled={busy !== null} onClick={() => void generate(round)} type="button">Generate draw</button>
                    )}
                    <button className="rounded-lg bg-slate-800 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50" disabled={busy !== null} onClick={() => void lock(round)} type="button">Lock round</button>
                  </div>
                )}

                {round.status === "OPEN" && round.pairingMethod === "MANUAL" && (
                  <div className="mt-5 grid gap-3 md:grid-cols-[1fr_1fr_auto]">
                  <select className="rounded-lg border border-slate-300 bg-white px-3 py-2" value={selectedTeams[`${round.id}:team1`] ?? ""} onChange={(event) => setSelectedTeams({ ...selectedTeams, [`${round.id}:team1`]: event.target.value })}>
                    <option value="">Team 1</option>
                    {availableTeams.map((team) => <option key={`one-${team.id}`} value={team.id}>{team.player1} / {team.player2}</option>)}
                  </select>
                  <select className="rounded-lg border border-slate-300 bg-white px-3 py-2" value={selectedTeams[`${round.id}:team2`] ?? ""} onChange={(event) => setSelectedTeams({ ...selectedTeams, [`${round.id}:team2`]: event.target.value })}>
                    <option value="">Team 2</option>
                    {availableTeams.map((team) => <option key={`two-${team.id}`} value={team.id}>{team.player1} / {team.player2}</option>)}
                  </select>
                  <button className="rounded-lg bg-amber-500 px-4 py-2 font-semibold text-white disabled:opacity-50" disabled={busy !== null || round.status !== "OPEN"} onClick={() => void createMatch(round)} type="button">Add match</button>
                  </div>
                )}

                {roundMatches.length > 0 && (
                  <div className="mt-5 border-t border-slate-200 pt-5">
                    <h3 className="text-sm font-bold uppercase tracking-wide text-slate-500">Matches</h3>
                    <div className="mt-3 space-y-2">
                      {roundMatches.map((match, index) => (
                        <div className="flex flex-wrap items-center justify-between gap-2 rounded-lg bg-slate-50 px-4 py-3 text-sm" key={match.id}>
                          <span className="font-semibold text-slate-800">
                            Match {index + 1}: {match.team1?.player1 ?? "TBD"} / {match.team1?.player2 ?? "TBD"} vs {match.team2?.player1 ?? "TBD"} / {match.team2?.player2 ?? "TBD"}
                          </span>
                          <span className="text-xs font-bold text-slate-500">{match.status}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                </section>
              );
            })}
          </div>
        </div>
      </main>
      <PublicFooter />
    </div>
  );
}

export default Draw;
