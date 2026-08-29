import { useCallback, useEffect, useMemo, useState } from "react";
import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { useParams } from "react-router-dom";
import { getEvent, listEvents } from "../api/eventApi";
import type { EventInfo } from "../types/event";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

interface Team {
  id: string;
  teamName: string;
  player1: string;
  player2: string;
}

interface Game {
  gameNumber: number;
  team1Score: number;
  team2Score: number;
}

interface Match {
  id: string;
  roundId: string;
  round: string;
  matchNumber: number;
  courtNumber: number | null;
  scheduledAt: string | null;
  status: string;
  team1: Team | null;
  team2: Team | null;
  games: Game[];
}

interface Round {
  id: string;
  roundNumber: number;
  roundName: string;
  status: string;
}

export default function LiveScores() {
  const navigate = useNavigate();
  const { eventId } = useParams();
  const [event, setEvent] = useState<EventInfo | null>(null);
  const [matches, setMatches] = useState<Match[]>([]);
  const [rounds, setRounds] = useState<Round[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadScores = useCallback(async () => {
    const resolvedEventId = eventId ?? (await listEvents()).find((item) => item.status !== "DRAFT")?.id ?? "";
    if (!resolvedEventId) {
      setEvent(null);
      setMatches([]);
      setRounds([]);
      return;
    }

    const [eventResponse, matchesResponse, roundsResponse] = await Promise.all([
      getEvent(resolvedEventId),
      fetch(`${API_BASE_URL}/api/v1/events/${resolvedEventId}/matches`),
      fetch(`${API_BASE_URL}/api/v1/events/${resolvedEventId}/rounds`),
    ]);

    if (!matchesResponse.ok || !roundsResponse.ok) {
      throw new Error("Unable to load live scores");
    }

    setEvent(eventResponse);
    setMatches((await matchesResponse.json()) ?? []);
    setRounds((await roundsResponse.json()) ?? []);
  }, [eventId]);

  useEffect(() => {
    void loadScores().catch((loadError) => {
      setError(loadError instanceof Error ? loadError.message : "Unable to load live scores");
    }).finally(() => setLoading(false));
  }, [loadScores]);

  const currentRound = rounds.find((round) => round.status === "IN_PROGRESS")
    ?? rounds.find((round) => round.status === "LOCKED")
    ?? rounds.find((round) => round.status === "OPEN")
    ?? rounds[rounds.length - 1];
  const currentMatches = currentRound
    ? matches.filter((match) => match.roundId === currentRound.id)
    : [];
  const liveMatches = matches.filter((match) => match.status === "IN_PROGRESS");
  const nextMatches = matches.filter((match) => match.status === "SCHEDULED").slice(0, 4);

  const displayScore = (match: Match, teamIndex: 0 | 1) => {
    const gameScores = match.games ?? [];
    if (gameScores.length === 0) return "-";
    return gameScores.map((game) => teamIndex === 0 ? game.team1Score : game.team2Score).join(" - ");
  };

  const statusLabel = useMemo(() => event?.status === "REGISTRATION_OPEN" ? "Registration open" : "Live", [event?.status]);

  if (loading) {
    return <PageShell><main className="flex flex-1 items-center justify-center p-10 text-slate-600">Loading live scores...</main></PageShell>;
  }

  if (error) {
    return <PageShell><main className="flex flex-1 items-center justify-center p-10 text-center"><div><h1 className="text-xl font-bold text-slate-950">Unable to load live scores</h1><p className="mt-2 text-slate-500">{error}</p></div></main></PageShell>;
  }

  return (
    <PageShell>
      <main className="flex-1 bg-slate-100 px-4 py-6 sm:py-10">
        <div className="mx-auto max-w-3xl">
          <section className="rounded-3xl bg-slate-950 p-6 text-white shadow-xl sm:p-8">
            <div className="flex items-start justify-between gap-5">
              <div>
                <p className="text-3xl" aria-hidden="true">🏸</p>
                <h1 className="mt-5 text-3xl font-black tracking-tight">ShuttleHub</h1>
                <p className="mt-1 text-sm font-semibold uppercase tracking-[0.18em] text-slate-400">Men's Doubles 2026</p>
              </div>
              <span className="inline-flex items-center gap-2 rounded-full bg-emerald-400/15 px-3 py-1.5 text-xs font-bold uppercase tracking-wide text-emerald-300"><span className="h-2 w-2 rounded-full bg-emerald-400" aria-hidden="true" />{statusLabel}</span>
            </div>
            <p className="mt-7 text-xs font-bold uppercase tracking-[0.2em] text-slate-400">Men's Doubles</p>
            <p className="mt-2 text-lg font-semibold text-white">{event?.name ?? "Tournament live scores"}</p>
          </section>

          <section className="mt-4 rounded-2xl bg-white p-5 shadow-sm sm:p-6">
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">🏆 Current round</p>
                <h2 className="mt-2 text-2xl font-black text-slate-950">{currentRound ? `Round ${currentRound.roundNumber}` : "Tournament setup"}</h2>
                <p className="mt-1 text-sm text-slate-500">{currentMatches.length} {currentMatches.length === 1 ? "match" : "matches"}</p>
              </div>
              <button className="rounded-lg bg-slate-950 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800" onClick={() => navigate("/live-scores#all-scores")} type="button">View Live Scores</button>
            </div>
          </section>

          <section className="mt-4 rounded-2xl bg-white p-5 shadow-sm sm:p-6">
            <div className="flex items-center gap-2"><span className="h-2.5 w-2.5 rounded-full bg-red-500" aria-hidden="true" /><h2 className="text-xs font-black uppercase tracking-[0.18em] text-slate-700">Live now</h2></div>
            <div className="mt-4 space-y-3">
              {liveMatches.length === 0 && <p className="rounded-xl bg-slate-50 p-4 text-sm text-slate-500">No matches are live right now.</p>}
              {liveMatches.map((match) => <ScoreRow key={match.id} match={match} score1={displayScore(match, 0)} score2={displayScore(match, 1)} live />)}
            </div>
          </section>

          <section className="mt-4 rounded-2xl bg-white p-5 shadow-sm sm:p-6">
            <h2 className="text-xs font-black uppercase tracking-[0.18em] text-slate-500">Next matches</h2>
            <div className="mt-4 divide-y divide-slate-200">
              {nextMatches.length === 0 && <p className="rounded-xl bg-slate-50 p-4 text-sm text-slate-500">No upcoming matches are scheduled.</p>}
              {nextMatches.map((match) => <div className="flex flex-wrap items-center justify-between gap-3 py-4 first:pt-0 last:pb-0" key={match.id}><div><p className="font-bold text-slate-800">{teamLabel(match.team1)} <span className="px-1 text-slate-400">vs</span> {teamLabel(match.team2)}</p><p className="mt-1 text-xs text-slate-500">Court {match.courtNumber ?? 1} <span className="px-1">•</span> {formatTime(match.scheduledAt)}</p></div><span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-bold text-slate-600">{match.round.replaceAll("_", " ")}</span></div>)}
            </div>
          </section>

          <nav className="mt-4 grid grid-cols-2 gap-3" aria-label="Score views">
            <a className="rounded-xl border border-slate-300 bg-white px-4 py-3 text-center text-sm font-black text-slate-800 transition hover:border-amber-400 hover:bg-amber-50" href="#bracket">Bracket</a>
            <a className="rounded-xl bg-amber-400 px-4 py-3 text-center text-sm font-black text-slate-950 transition hover:bg-amber-300" href="#all-scores">All scores</a>
          </nav>

          <section className="mt-4 rounded-2xl bg-white p-5 shadow-sm sm:p-6" id="all-scores">
            <h2 className="text-xs font-black uppercase tracking-[0.18em] text-slate-500">All scores</h2>
            <div className="mt-4 space-y-2">
              {matches.filter((match) => match.status === "COMPLETED").map((match) => <ScoreRow key={match.id} match={match} score1={displayScore(match, 0)} score2={displayScore(match, 1)} />)}
              {matches.filter((match) => match.status === "COMPLETED").length === 0 && <p className="text-sm text-slate-500">Completed scores will appear here.</p>}
            </div>
          </section>

          <section className="mt-4 rounded-2xl bg-white p-5 shadow-sm sm:p-6" id="bracket">
            <h2 className="text-xs font-black uppercase tracking-[0.18em] text-slate-500">Bracket</h2>
            <div className="mt-4 flex flex-wrap gap-2">{rounds.map((round) => <span className="rounded-full bg-slate-100 px-3 py-2 text-xs font-bold uppercase tracking-wide text-slate-600" key={round.id}>{round.roundName.replaceAll("_", " ")}</span>)}</div>
          </section>
        </div>
      </main>
    </PageShell>
  );
}

function ScoreRow({ match, score1, score2, live = false }: { match: Match; score1: number | string; score2: number | string; live?: boolean }) {
  return <div className="rounded-xl border border-slate-200 bg-slate-50 p-4"><div className="mb-3 flex items-center justify-between"><span className="text-xs font-black uppercase tracking-wide text-slate-400">Court {match.courtNumber ?? 1}</span>{live ? <span className="text-xs font-black uppercase text-red-600">Live</span> : <span className="text-xs font-bold text-slate-500">{match.status}</span>}</div><div className="grid grid-cols-[1fr_auto] gap-x-4 gap-y-2 text-sm"><span className="font-bold text-slate-800">{teamLabel(match.team1)}</span><span className="text-base font-black text-slate-950">{score1}</span><span className="font-bold text-slate-800">{teamLabel(match.team2)}</span><span className="text-base font-black text-slate-950">{score2}</span></div></div>;
}

function teamLabel(team: Team | null) {
  return team ? `${team.player1} - ${team.player2}` : "TBD";
}

function formatTime(value: string | null) {
  if (!value) return "Time TBD";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "Time TBD" : date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
}

function PageShell({ children }: { children: ReactNode }) {
  return <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900"><PublicHeader />{children}<PublicFooter /></div>;
}
