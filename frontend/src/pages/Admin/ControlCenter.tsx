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
  status: string;
}

interface Team {
  id: string;
  teamName?: string;
  player1: string;
  player2: string;
}

interface Match {
  roundId: string;
  round: string;
  status: string;
  team1: Team | null;
  team2: Team | null;
  winnerTeamId: string | null;
  loserTeamId: string | null;
}

export default function ControlCenter() {
  const navigate = useNavigate();
  const [rounds, setRounds] = useState<Round[]>([]);
  const [matches, setMatches] = useState<Match[]>([]);
  const [teamCount, setTeamCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const request = useCallback(async (url: string) => {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      headers: {
        Authorization: `Bearer ${getAdminToken()}`,
      },
    });

    if (response.status === 401 || response.status === 403) {
      logout();
      navigate("/admin/login", { replace: true });
      throw new Error("Session expired");
    }

    if (!response.ok) {
      throw new Error("Unable to load tournament summary");
    }

    return response.json();
  }, [navigate]);

  useEffect(() => {
    async function loadSummary() {
      try {
        const [roundData, matchData, teamData] = await Promise.all([
          request(`/api/v1/admin/events/${EVENT_ID}/rounds`),
          request(`/api/v1/events/${EVENT_ID}/matches`),
          request(`/api/v1/events/${EVENT_ID}/teams`),
        ]);
        setRounds(roundData ?? []);
        setMatches(matchData ?? []);
        setTeamCount((teamData?.teams ?? teamData ?? []).length);
      } catch (loadError) {
        if (loadError instanceof Error && loadError.message !== "Session expired") {
          setError(loadError.message);
        }
      } finally {
        setLoading(false);
      }
    }

    void loadSummary();
  }, [request]);

  const isRoundComplete = (round: Round) => {
    const roundMatches = matches.filter((match) => match.roundId === round.id);
    return round.status === "COMPLETED"
      || (roundMatches.length > 0 && roundMatches.every((match) => (
        match.status === "COMPLETED" || match.status === "CANCELLED"
      )));
  };
  const currentRound = rounds.find((round) => !isRoundComplete(round)) ?? rounds[rounds.length - 1];
  const currentMatches = currentRound ? matches.filter((match) => match.roundId === currentRound.id) : [];
  const completedMatches = matches.filter((match) => match.status === "COMPLETED").length;
  const finalMatch = matches.find((match) => match.round === "FINAL" && match.status === "COMPLETED");
  const finalWinner = finalMatch?.winnerTeamId
    ? [finalMatch.team1, finalMatch.team2].find((team) => team?.id === finalMatch.winnerTeamId) ?? null
    : null;
  const finalRunnerUp = finalMatch?.loserTeamId
    ? [finalMatch.team1, finalMatch.team2].find((team) => team?.id === finalMatch.loserTeamId) ?? null
    : null;
  const tournamentComplete = Boolean(finalWinner && finalRunnerUp);
  const stageLabel = (round: Round) => {
    if (round.roundName.startsWith("SEMIFINAL")) {
      return "Semifinal";
    }
    if (round.roundName === "FINAL") {
      return "Final";
    }
    return `R${round.roundNumber}`;
  };
  const progress = [
    { label: "Registration", complete: teamCount > 0 },
    ...rounds
      .filter((round, index, allRounds) => (
        !round.roundName.startsWith("SEMIFINAL")
        || allRounds.findIndex((candidate) => candidate.roundName.startsWith("SEMIFINAL")) === index
      ))
      .map((round) => ({
        label: stageLabel(round),
        complete: isRoundComplete(round) || (round.roundName === "FINAL" && tournamentComplete),
      })),
  ];

  if (loading) {
    return (
      <div className="flex min-h-screen flex-col bg-slate-100">
        <PublicHeader />
        <main className="flex flex-1 items-center justify-center p-10 text-slate-700">Loading control center...</main>
        <PublicFooter />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900">
      <PublicHeader />
      <main className="flex-1 px-4 py-8">
        <div className="mx-auto max-w-5xl">
          <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.18em] text-slate-500">ShuttleHub · Admin</p>
              <h1 className="mt-2 text-3xl font-black text-slate-950">Admin Control Center</h1>
              <p className="mt-2 text-sm text-slate-500">Men&apos;s Doubles 2026</p>
            </div>
            <button className="rounded-lg border border-red-200 bg-red-50 px-4 py-2 text-sm font-bold text-red-700 transition hover:bg-red-100" onClick={() => { logout(); navigate("/admin/login", { replace: true }); }} type="button">Sign out</button>
          </div>

          {error && <p className="mb-5 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm font-medium text-amber-900" role="alert">{error}</p>}

          <section className="mb-6 rounded-2xl bg-white p-6 shadow-sm">
            <div className="flex flex-wrap items-end justify-between gap-4">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Men&apos;s Doubles</p>
                <h2 className="mt-2 text-2xl font-bold text-slate-950">Tournament overview</h2>
              </div>
              <span className="inline-flex items-center gap-2 rounded-full bg-emerald-100 px-3 py-1.5 text-xs font-bold uppercase tracking-wide text-emerald-800"><span className="h-2 w-2 rounded-full bg-emerald-500" aria-hidden="true" />{tournamentComplete ? "Complete" : "Live"}</span>
            </div>
            <div className="mt-6 grid gap-4 sm:grid-cols-3">
              <SummaryTile label="Teams" value={teamCount} />
              <SummaryTile label="Matches" value={matches.length} />
              <SummaryTile label="Current round" value={currentRound ? `Round ${currentRound.roundNumber}` : "-"} />
            </div>
          </section>

          <section className="mb-6 rounded-2xl bg-slate-950 p-6 text-white shadow-sm">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">Tournament progress</p>
                <h2 className="mt-2 text-xl font-bold">Road to the final</h2>
              </div>
              <span className="text-sm font-semibold text-slate-400">{completedMatches} matches completed</span>
            </div>
            <div className="mt-6 flex flex-wrap items-center gap-y-3">
              {progress.map((stage, index) => (
                <div className="flex items-center" key={`${stage.label}-${index}`}>
                  <span className={`rounded-full px-3 py-2 text-xs font-bold ${stage.complete ? "bg-emerald-400 text-emerald-950" : currentRound && stage.label === stageLabel(currentRound) ? "bg-amber-300 text-amber-950" : "bg-white/10 text-slate-400"}`}>{stage.complete ? "✓" : "○"} {stage.label}</span>
                  {index < progress.length - 1 && <span className="px-2 text-slate-600" aria-hidden="true">→</span>}
                </div>
              ))}
            </div>
          </section>

          <section className="mb-6 rounded-2xl bg-white p-6 shadow-sm">
            <div className="flex flex-wrap items-end justify-between gap-3">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Current round</p>
                <h2 className="mt-2 text-2xl font-bold text-slate-950">{tournamentComplete ? "Tournament complete" : currentRound ? `Round ${currentRound.roundNumber}` : "No round started"}</h2>
                <p className="mt-1 text-sm text-slate-500">{tournamentComplete ? "Final result" : `${currentMatches.length} ${currentMatches.length === 1 ? "match" : "matches"} · ${currentRound?.status ?? "Awaiting setup"}`}</p>
              </div>
              <button className="rounded-lg bg-slate-950 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800" onClick={() => navigate("/admin/draw")} type="button">Manage round</button>
            </div>
            <div className="mt-5 space-y-3">
              {tournamentComplete && finalWinner && finalRunnerUp && (
                <div className="grid gap-4 rounded-xl border border-amber-200 bg-amber-50 p-5 sm:grid-cols-2">
                  <div>
                    <p className="text-xs font-bold uppercase tracking-[0.16em] text-amber-700">Champion</p>
                    <p className="mt-2 text-lg font-black text-slate-950">{finalWinner.player1} / {finalWinner.player2}</p>
                    {finalWinner.teamName && <p className="mt-1 text-sm text-slate-600">{finalWinner.teamName}</p>}
                  </div>
                  <div>
                    <p className="text-xs font-bold uppercase tracking-[0.16em] text-slate-500">Runner-up</p>
                    <p className="mt-2 text-lg font-black text-slate-950">{finalRunnerUp.player1} / {finalRunnerUp.player2}</p>
                    {finalRunnerUp.teamName && <p className="mt-1 text-sm text-slate-600">{finalRunnerUp.teamName}</p>}
                  </div>
                </div>
              )}
              {!tournamentComplete && currentMatches.slice(0, 3).map((match, index) => (
                <div className="flex flex-col gap-2 rounded-xl border border-slate-200 bg-slate-50 p-4 sm:flex-row sm:items-center sm:justify-between" key={`${match.roundId}-${index}`}>
                  <span className="font-black text-slate-400">M{index + 1}</span>
                  <span className="flex-1 px-3 text-sm font-semibold text-slate-800">{match.team1?.player1 ?? "TBD"} / {match.team1?.player2 ?? "TBD"} <span className="px-2 text-slate-400">VS</span> {match.team2?.player1 ?? "TBD"} / {match.team2?.player2 ?? "TBD"}</span>
                  <span className={`text-xs font-bold ${match.status === "COMPLETED" ? "text-emerald-700" : "text-amber-700"}`}>{match.status}</span>
                </div>
              ))}
              {currentMatches.length > 3 && <p className="pt-1 text-center text-sm font-semibold text-slate-500">+ {currentMatches.length - 3} more matches</p>}
              {currentMatches.length === 0 && <p className="rounded-lg bg-slate-50 p-4 text-sm text-slate-500">No matches are available yet.</p>}
            </div>
          </section>

          <section aria-labelledby="quick-actions-heading">
            <h2 className="mb-3 text-xs font-bold uppercase tracking-[0.18em] text-slate-500" id="quick-actions-heading">Quick actions</h2>
            <div className="grid gap-4 sm:grid-cols-3">
              <ActionTile label="Registrations" description="Confirm teams and manage entries" onClick={() => navigate("/admin/registrations")} />
              <ActionTile label="Draw & results" description="Create rounds and enter scores" onClick={() => navigate("/admin/draw")} />
              <ActionTile label="Live scores" description="View the public tournament board" onClick={() => navigate("/")} />
            </div>
          </section>
        </div>
      </main>
      <PublicFooter />
    </div>
  );
}

function SummaryTile({ label, value }: { label: string; value: string | number }) {
  return <div className="rounded-xl border border-slate-200 bg-slate-50 p-5"><p className="text-3xl font-black text-slate-950">{value}</p><p className="mt-2 text-xs font-bold uppercase tracking-wide text-slate-500">{label}</p></div>;
}

function ActionTile({ label, description, onClick }: { label: string; description: string; onClick: () => void }) {
  return <button className="rounded-xl border border-slate-200 bg-white p-5 text-left shadow-sm transition hover:-translate-y-0.5 hover:border-amber-300 hover:shadow-md" onClick={onClick} type="button"><p className="font-bold text-slate-950">{label} <span aria-hidden="true">→</span></p><p className="mt-2 text-sm text-slate-500">{description}</p></button>;
}
