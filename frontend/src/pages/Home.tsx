import { useEffect, useState } from "react";
import { listEvents } from "../api/eventApi";
import type { EventInfo } from "../types/event";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

function Home() {
  const [events, setEvents] = useState<EventInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    void listEvents()
      .then(setEvents)
      .catch((loadError: unknown) => {
        setError(loadError instanceof Error ? loadError.message : "Failed to load events");
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <PublicHeader />
      <main>
        <section className="border-b border-white/10 px-5 pb-16 pt-14 sm:px-8 sm:pt-20">
          <div className="mx-auto max-w-6xl">
            <div className="max-w-3xl">
              <p className="text-sm font-bold uppercase tracking-[0.2em] text-amber-300">ShuttleHub tournament platform</p>
              <h1 className="mt-5 text-4xl font-black tracking-tight sm:text-6xl">Every match, clearly organized.</h1>
              <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-300">Discover upcoming competitions, register teams, follow live scores, and keep every round moving from first draw to final point.</p>
              <div className="mt-8 flex flex-wrap gap-3">
                <a className="rounded-xl bg-amber-400 px-5 py-3 font-bold text-slate-950 transition hover:bg-amber-300" href="#events">Explore events</a>
                <a className="rounded-xl border border-slate-700 bg-slate-900 px-5 py-3 font-semibold text-white transition hover:bg-slate-800" href="#benefits">Why ShuttleHub</a>
              </div>
            </div>
          </div>
        </section>

        <section className="bg-slate-100 px-5 py-12 text-slate-900 sm:px-8" id="events">
          <div className="mx-auto max-w-6xl">
            <div className="flex flex-wrap items-end justify-between gap-4">
              <div>
                <p className="text-sm font-bold uppercase tracking-[0.2em] text-slate-500">Tournament calendar</p>
                <h2 className="mt-2 text-3xl font-black text-slate-950">Upcoming events</h2>
              </div>
              <span className="text-sm font-semibold text-slate-500">{events.length} {events.length === 1 ? "event" : "events"}</span>
            </div>
            {loading && <p className="mt-6 rounded-2xl bg-white p-8 text-center text-slate-500 shadow-sm">Loading events...</p>}
            {error && <p className="mt-6 rounded-2xl border border-amber-200 bg-amber-50 p-5 text-amber-900" role="alert">{error}</p>}
            {!loading && !error && events.length === 0 && <p className="mt-6 rounded-2xl bg-white p-8 text-center text-slate-500 shadow-sm">No published events are available yet.</p>}
            {!loading && !error && events.length > 0 && <div className="mt-6 grid gap-5 md:grid-cols-2 lg:grid-cols-3">{events.map((event) => <EventCard event={event} key={event.id} />)}</div>}
          </div>
        </section>

        <section className="px-5 py-14 sm:px-8" id="benefits">
          <div className="mx-auto max-w-6xl">
            <div className="max-w-2xl"><p className="text-sm font-bold uppercase tracking-[0.2em] text-amber-300">Built for tournament day</p><h2 className="mt-3 text-3xl font-black">A smoother experience for players and organizers.</h2></div>
            <div className="mt-8 grid gap-5 md:grid-cols-3"><Feature title="Simple registration" text="Collect complete team details in one focused flow, with capacity and registration status always visible." /><Feature title="Confident operations" text="Create draws, manage rounds, close registration, and enter results from one reliable admin workflow." /><Feature title="Clear live tracking" text="Give players and spectators a shared view of courts, schedules, scores, and tournament progress." /></div>
          </div>
        </section>
      </main>
      <PublicFooter />
    </div>
  );
}

function EventCard({ event }: { event: EventInfo }) {
  const open = event.status === "REGISTRATION_OPEN";
  return <article className="flex flex-col rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"><div className="flex items-start justify-between gap-3"><div><p className="text-xs font-bold uppercase tracking-[0.16em] text-slate-500">Men&apos;s Doubles</p><h3 className="mt-2 text-xl font-black text-slate-950">{event.name}</h3></div><span className={`rounded-full px-3 py-1 text-xs font-bold ${open ? "bg-emerald-100 text-emerald-800" : "bg-slate-100 text-slate-600"}`}>{open ? "Registration open" : event.status.replaceAll("_", " ")}</span></div><div className="mt-6 grid grid-cols-2 gap-4 border-y border-slate-200 py-4"><Stat label="Teams" value={event.maxTeams === null ? `${event.registeredTeams}` : `${event.registeredTeams} / ${event.maxTeams}`} /><Stat label="Format" value="Doubles" /></div><div className="mt-auto flex flex-wrap gap-2 pt-5"><a className="rounded-lg bg-slate-950 px-3 py-2 text-sm font-bold text-white transition hover:bg-slate-800" href={`/events/${event.id}`}>View event</a>{open && <a className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-bold text-slate-700 transition hover:bg-slate-50" href={`/events/${event.id}/register`}>Register</a>}</div></article>;
}

function Stat({ label, value }: { label: string; value: string }) {
  return <div><p className="text-xs font-bold uppercase tracking-wide text-slate-400">{label}</p><p className="mt-1 font-bold text-slate-800">{value}</p></div>;
}

function Feature({ title, text }: { title: string; text: string }) {
  return <article className="rounded-2xl border border-slate-800 bg-slate-900 p-6"><h3 className="text-lg font-bold">{title}</h3><p className="mt-3 text-sm leading-7 text-slate-300">{text}</p></article>;
}

export default Home;
