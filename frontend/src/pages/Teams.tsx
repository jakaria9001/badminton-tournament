import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { getTeams } from "../api/teamApi";
import type { Team } from "../api/teamApi";

import { getEvent, listEvents } from "../api/eventApi";
import type { EventInfo } from "../types/event";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

function Teams() {
  const { eventId } = useParams();
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [event, setEvent] = useState<EventInfo | null>(null);

  useEffect(() => {
    async function loadTeams() {
      try {
        const resolvedEventId = eventId ?? (await listEvents()).find((item) => item.status !== "DRAFT")?.id ?? "";
        if (!resolvedEventId) {
          setTeams([]);
          setEvent(null);
          return;
        }

        const [teamsData, eventData] = await Promise.all([
          getTeams(resolvedEventId),
          getEvent(resolvedEventId),
        ]);

        setTeams(teamsData);
        setEvent(eventData);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Failed to load teams",
        );
      } finally {
        setLoading(false);
      }
    }

    void loadTeams();
  }, [eventId]);

  return (
    <div className="flex min-h-screen flex-col bg-slate-950 text-white">
      <PublicHeader />
      <main className="flex-1 px-4 py-8">
      <div className="mx-auto max-w-4xl">
        <div className="mb-8 rounded-[28px] border border-slate-800 bg-slate-950 p-8 shadow-2xl shadow-slate-950/40">
          <div className="flex flex-col items-center text-center">
            <h1 className="mt-4 text-3xl font-black tracking-tight sm:text-4xl">
              Confirmed Teams
            </h1>

            <p className="mt-2 text-slate-300">
              Men's Doubles 2026 tournament lineup
            </p>
          </div>
        </div>

        {loading && (
          <div className="rounded-2xl border border-slate-800 bg-slate-900 p-10 text-center text-slate-300">
            Loading teams...
          </div>
        )}

        {error && (
          <div className="rounded-2xl border border-red-900 bg-red-950/40 p-4 text-red-200">
            {error}
          </div>
        )}

        {!loading && !error && (
          <>
            <div className="mb-5 flex items-center justify-between rounded-2xl border border-slate-800 bg-slate-900/80 p-4">
              <h2 className="text-xl font-semibold text-white">
                Confirmed Teams
              </h2>

              <span className="rounded-full border border-slate-700 bg-slate-800 px-3 py-1 text-sm text-slate-200">
                {event?.maxTeams !== null && event?.maxTeams !== undefined
                  ? `${teams.length} / ${event.maxTeams} Teams`
                  : `${teams.length} Teams`}
              </span>
            </div>

            <div className="overflow-hidden rounded-2xl border border-slate-800 bg-slate-900/80 shadow-lg shadow-slate-950/20">
              <div className="hidden grid-cols-[56px_1.15fr_1.25fr_100px] gap-3 border-b border-slate-800 bg-slate-950/60 px-4 py-3 text-xs font-semibold uppercase tracking-[0.12em] text-slate-400 md:grid">
                <span>Rank</span>
                <span>Team</span>
                <span>Players</span>
                <span className="text-right">Status</span>
              </div>

              {teams.map((team, index) => (
                <div
                  key={team.id}
                  className="grid gap-3 border-b border-slate-800 px-4 py-3 last:border-b-0 md:grid-cols-[56px_1.15fr_1.25fr_100px] md:items-center"
                >
                  <div className="flex items-center">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white text-sm font-bold text-slate-950">
                      {index + 1}
                    </div>
                  </div>

                  <div>
                    <h3 className="font-semibold text-white">{team.teamName}</h3>
                  </div>

                  <div className="text-sm text-slate-300">
                    {team.player1.name}
                    {" / "}
                    {team.player2.name}
                  </div>

                  <div className="md:text-right">
                    <span className="inline-flex rounded-full border border-amber-400/30 bg-amber-500/10 px-2.5 py-1 text-xs font-medium text-amber-200">
                      {team.status}
                    </span>
                  </div>
                </div>
              ))}
            </div>

            {teams.length === 0 && (
              <div className="rounded-2xl border border-slate-800 bg-slate-900/80 p-10 text-center text-slate-300">
                No teams registered yet.
              </div>
            )}
          </>
        )}
      </div>
      </main>
      <PublicFooter />
    </div>
  );
}

export default Teams;