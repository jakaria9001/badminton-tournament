import { useEffect, useState } from "react";
import { getEvent } from "../api/eventApi";
import type { EventInfo } from "../types/event";

const EVENT_ID =
  "00000000-0000-0000-0000-000000000002";

function Home() {
  const [event, setEvent] =
    useState<EventInfo | null>(null);

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");

  useEffect(() => {
    async function loadEvent() {
      try {
        const data = await getEvent(EVENT_ID);
        setEvent(data);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Failed to load tournament",
        );
      } finally {
        setLoading(false);
      }
    }

    loadEvent();
  }, []);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 text-white">
        Loading tournament...
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 text-center text-white">
        <div>
          <p className="text-xl font-semibold">
            Unable to load tournament
          </p>

          <p className="mt-2 text-slate-400">
            {error}
          </p>
        </div>
      </div>
    );
  }

  const isOpen =
    event.status === "REGISTRATION_OPEN";

  const capacityText =
    event.maxTeams !== null
      ? `${event.registeredTeams} / ${event.maxTeams}`
      : `${event.registeredTeams}`;

  const percentage =
    event.maxTeams
      ? Math.min(
          100,
          (event.registeredTeams /
            event.maxTeams) *
            100,
        )
      : 0;

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-sm">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-5 py-4 sm:px-8">
          <a href="/" className="flex items-center gap-3" aria-label="ShuttleHub home">
            <img
              src="/shuttlehub_logo.svg"
              alt="ShuttleHub logo"
              className="h-10 w-auto"
            />
            <div>
              <div className="text-lg font-bold tracking-tight text-white">ShuttleHub</div>
              <div className="text-[10px] uppercase tracking-[0.22em] text-slate-400">Men's Doubles</div>
            </div>
          </a>

          <nav className="hidden items-center gap-6 text-sm text-slate-300 md:flex">
            <a href="/" className="transition hover:text-white">Home</a>
            <a href="/teams" className="transition hover:text-white">Teams</a>
            <a href="/register" className="transition hover:text-white">Register</a>
          </nav>
        </div>
      </header>

      <section className="px-5 pb-16 pt-10 sm:px-8">
        <div className="mx-auto max-w-6xl">
          <div className="overflow-hidden rounded-[32px] border border-slate-800 bg-gradient-to-br from-slate-900 via-slate-900 to-slate-950 shadow-2xl shadow-slate-950/40">
            <div className="grid gap-8 px-6 py-8 md:grid-cols-[1.2fr_0.8fr] md:px-10 md:py-10">
              <div className="flex flex-col justify-center">
                <div className="mb-5 inline-flex w-fit items-center gap-2 rounded-full border border-amber-400/30 bg-amber-500/10 px-4 py-2 text-sm font-medium text-amber-200">
                  <span>🏸</span>
                  2026 Tournament Series
                </div>

                <h1 className="text-4xl font-black tracking-tight sm:text-5xl lg:text-7xl">
                  <a href="/" className="transition hover:text-slate-300">ShuttleHub</a>
                  <span className="mt-2 block text-slate-400">
                    Men's Doubles 2026
                  </span>
                </h1>

                <p className="mt-5 max-w-xl text-base leading-8 text-slate-300 sm:text-lg">
                  A competitive, club-style badminton experience built for strong matchups, team spirit, and a polished tournament day.
                </p>

                <div className="mt-8 flex flex-wrap gap-3">
                  {isOpen && (
                    <button
                      className="rounded-xl bg-white px-6 py-3 font-bold text-slate-950 transition hover:bg-slate-200"
                      onClick={() => {
                        window.location.href = "/register";
                      }}
                    >
                      Register Your Team
                    </button>
                  )}

                  <button
                    className="rounded-xl border border-slate-700 bg-slate-800 px-6 py-3 font-semibold text-white transition hover:bg-slate-700"
                    onClick={() => {
                      window.location.href = "/teams";
                    }}
                  >
                    View Teams
                  </button>
                </div>

                <div className="mt-8 flex flex-wrap gap-4 text-sm text-slate-300">
                  <div className="rounded-full border border-slate-700 bg-slate-800 px-3 py-2">
                    1-night showdown
                  </div>
                  <div className="rounded-full border border-slate-700 bg-slate-800 px-3 py-2">
                    Competitive spirit
                  </div>
                  <div className="rounded-full border border-slate-700 bg-slate-800 px-3 py-2">
                    Elite club experience
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-center">
                <div className="w-full max-w-md rounded-[28px] border border-slate-800 bg-slate-950 p-5 shadow-xl shadow-slate-950/30">
                  <div className="mb-6 flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">Event</p>
                      <h2 className="mt-1 text-2xl font-bold text-white">{event.name}</h2>
                    </div>
                    <div className="rounded-full border border-emerald-700 bg-emerald-950/70 px-3 py-1 text-sm font-semibold text-emerald-300">
                      {isOpen ? "Open" : "Closed"}
                    </div>
                  </div>

                  <div className="space-y-4">
                    <InfoCard label="Venue" value="Elite Shuttler Club" />
                    <InfoCard label="Category" value="Men's Doubles" />
                    <InfoCard label="Teams" value={capacityText} />
                  </div>

                  <div className="mt-5 rounded-2xl border border-slate-800 bg-slate-900 p-4">
                    <div className="flex items-center justify-between text-sm text-slate-300">
                      <span>Capacity</span>
                      <span>{event.maxTeams !== null ? `${Math.min(100, Math.round(percentage))}%` : "N/A"}</span>
                    </div>
                    {event.maxTeams !== null && (
                      <div className="mt-3 h-2.5 overflow-hidden rounded-full bg-slate-800">
                        <div
                          className="h-full rounded-full bg-gradient-to-r from-amber-300 via-amber-400 to-amber-500 transition-all"
                          style={{ width: `${percentage}%` }}
                        />
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="border-y border-white/10 bg-white/5 px-5 py-8">
        <div className="mx-auto grid max-w-5xl gap-4 sm:grid-cols-3">
          <InfoCard label="Event" value={event.name} />
          <InfoCard label="Registration" value={isOpen ? "Open" : "Closed"} />
          <InfoCard label="Teams" value={capacityText} />
        </div>
      </section>

      <section className="px-5 py-12">
        <div className="mx-auto max-w-5xl">
          <div className="mb-6 text-center">
            <p className="text-sm font-medium uppercase tracking-[0.2em] text-slate-400">
              Tournament Snapshot
            </p>
            <h2 className="mt-3 text-3xl font-bold text-white">
              Build momentum before the final whistle
            </h2>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            <FeatureCard
              title="Fast registration"
              description="Simple two-player sign-up with quick team submission and instant confirmation feedback."
              icon="⚡"
            />
            <FeatureCard
              title="Live standings"
              description="Track the event as teams fill up and competitors rise through the bracket."
              icon="📊"
            />
            <FeatureCard
              title="Club atmosphere"
              description="A competitive but welcoming environment designed for fun, skill, and community."
              icon="🏆"
            />
          </div>
        </div>
      </section>

      <section className="px-5 pb-12">
        <div className="mx-auto max-w-5xl rounded-2xl border border-white/10 bg-white/5 p-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <p className="text-sm text-slate-400">Tournament capacity</p>
              <p className="mt-1 text-xl font-bold text-white">
                {capacityText} teams registered
              </p>
            </div>

            <div className="flex items-center gap-3 rounded-full border border-white/10 bg-slate-950/60 px-3 py-2 text-sm text-slate-300">
              <span className="text-xl">🏸</span>
              {event.maxTeams !== null && event.maxTeams - event.registeredTeams > 0
                ? `${event.maxTeams - event.registeredTeams} spots remaining`
                : "Tournament full"}
            </div>
          </div>

          {event.maxTeams !== null && (
            <div className="mt-5">
              <div className="h-2.5 overflow-hidden rounded-full bg-slate-800">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-amber-300 via-amber-400 to-amber-500 transition-all"
                  style={{ width: `${percentage}%` }}
                />
              </div>
            </div>
          )}
        </div>
      </section>

      <footer className="border-t border-white/10 px-5 py-8">
        <div className="mx-auto max-w-5xl text-center text-sm text-slate-500">
          Elite Shuttler Club &copy; 2026. All rights reserved.
        </div>
      </footer>
    </div>
  );
}

interface InfoCardProps {
  label: string;
  value: string;
}

function InfoCard({ label, value }: InfoCardProps) {
  return (
    <div className="rounded-xl border border-white/10 bg-slate-950/80 p-5 shadow-lg shadow-slate-950/20">
      <p className="text-sm text-slate-500">{label}</p>
      <p className="mt-2 font-semibold text-white">{value}</p>
    </div>
  );
}

interface FeatureCardProps {
  title: string;
  description: string;
  icon: string;
}

function FeatureCard({ title, description, icon }: FeatureCardProps) {
  return (
    <div className="rounded-2xl border border-white/10 bg-slate-950/70 p-5 shadow-lg shadow-slate-950/20">
      <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 text-2xl">
        {icon}
      </div>
      <h3 className="mt-4 text-xl font-bold text-white">{title}</h3>
      <p className="mt-2 text-sm leading-6 text-slate-300">{description}</p>
    </div>
  );
}

export default Home;