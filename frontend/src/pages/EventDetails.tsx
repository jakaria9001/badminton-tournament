import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getEvent } from "../api/eventApi";
import { getAdminProfile, getAdminToken, type AdminProfile } from "../api/authApi";
import type { EventInfo } from "../types/event";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

export default function EventDetails() {
  const { eventId } = useParams();
  const navigate = useNavigate();
  const [event, setEvent] = useState<EventInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [profile, setProfile] = useState<AdminProfile | null>(null);

  useEffect(() => {
    if (!eventId) {
      setError("Event not found");
      setLoading(false);
      return;
    }

    void getEvent(eventId)
      .then(setEvent)
      .catch((loadError: unknown) => {
        setError(loadError instanceof Error ? loadError.message : "Failed to load event");
      })
      .finally(() => setLoading(false));
  }, [eventId]);

  useEffect(() => {
    if (!getAdminToken()) {
      return;
    }

    void getAdminProfile().then(setProfile).catch(() => undefined);
  }, []);

  if (loading) {
    return <PageShell><main className="flex flex-1 items-center justify-center p-10 text-slate-600">Loading event...</main></PageShell>;
  }

  if (error || !event) {
    return <PageShell><main className="flex flex-1 items-center justify-center p-10 text-center"><div><h1 className="text-xl font-bold text-slate-950">Unable to load event</h1><p className="mt-2 text-slate-500">{error}</p><button className="mt-5 rounded-lg bg-slate-950 px-4 py-2 text-sm font-bold text-white" onClick={() => navigate("/")} type="button">Back to events</button></div></main></PageShell>;
  }

  const open = event.status === "REGISTRATION_OPEN";
  const capacity = event.maxTeams === null ? `${event.registeredTeams} teams registered` : `${event.registeredTeams} / ${event.maxTeams} teams`;
  const percentage = event.maxTeams ? Math.min(100, (event.registeredTeams / event.maxTeams) * 100) : 0;
  const canManageRegistrations = profile?.role === "SUPER_ADMIN" || profile?.eventId === event.id;

  return <PageShell><main className="flex-1 bg-slate-100 px-4 py-8 sm:py-12"><div className="mx-auto max-w-5xl"><button className="mb-5 text-sm font-bold text-slate-600 transition hover:text-slate-950" onClick={() => navigate("/")} type="button">← All events</button><section className="overflow-hidden rounded-3xl bg-slate-950 p-7 text-white shadow-xl sm:p-10"><p className="text-sm font-bold uppercase tracking-[0.2em] text-amber-300">ShuttleHub tournament series</p><div className="mt-4 flex flex-wrap items-start justify-between gap-5"><div><h1 className="text-4xl font-black tracking-tight sm:text-5xl">{event.name}</h1><p className="mt-4 max-w-2xl text-base leading-7 text-slate-300">A professionally managed men's doubles tournament with organized draws, transparent results, and live match tracking from the opening game to the final.</p></div><span className="rounded-full bg-emerald-400/15 px-3 py-1.5 text-xs font-bold uppercase tracking-wide text-emerald-300">{open ? "Registration open" : event.status.replaceAll("_", " ")}</span></div></section><section className="mt-5 grid gap-4 sm:grid-cols-3"><Detail label="Venue" value="Elite Shuttler Club" /><Detail label="Category" value="Men's doubles" /><Detail label="Teams" value={capacity} /></section><section className="mt-5 rounded-2xl bg-white p-6 shadow-sm sm:p-8"><div className="flex flex-wrap items-end justify-between gap-3"><div><p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Registration status</p><h2 className="mt-2 text-2xl font-black text-slate-950">{open ? "Build your team for tournament day" : "Registration is currently closed"}</h2></div><span className="text-sm font-bold text-slate-500">{event.maxTeams === null ? `${event.registeredTeams} teams` : `${Math.round(percentage)}% full`}</span></div>{event.maxTeams !== null && <div className="mt-5 h-3 overflow-hidden rounded-full bg-slate-200"><div className="h-full rounded-full bg-amber-400" style={{ width: `${percentage}%` }} /></div>}<p className="mt-4 max-w-2xl text-sm leading-6 text-slate-500">{open ? "Reserve a place for your pair, then follow the draw and scores as the event progresses." : "You can still browse confirmed teams and follow live scores for this event."}</p><div className="mt-6 flex flex-wrap gap-3">{open && <button className="rounded-lg bg-amber-400 px-4 py-2.5 text-sm font-bold text-slate-950 transition hover:bg-amber-300" onClick={() => navigate(`/events/${event.id}/register`)} type="button">Register team</button>}<button className="rounded-lg bg-slate-950 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800" onClick={() => navigate(`/events/${event.id}/live-scores`)} type="button">Live scores</button><button className="rounded-lg border border-slate-300 px-4 py-2.5 text-sm font-bold text-slate-700 transition hover:bg-slate-50" onClick={() => navigate(`/events/${event.id}/teams`)} type="button">View teams</button>{canManageRegistrations && <button className="rounded-lg border border-amber-400 bg-amber-50 px-4 py-2.5 text-sm font-bold text-amber-900 transition hover:bg-amber-100" onClick={() => navigate(`/admin/events/${event.id}/registrations`)} type="button">Manage registrations</button>}</div></section></div></main></PageShell>;
}

function Detail({ label, value }: { label: string; value: string }) {
  return <div className="rounded-2xl bg-white p-5 shadow-sm"><p className="text-xs font-bold uppercase tracking-wide text-slate-400">{label}</p><p className="mt-2 font-bold text-slate-800">{value}</p></div>;
}

function PageShell({ children }: { children: ReactNode }) {
  return <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900"><PublicHeader />{children}<PublicFooter /></div>;
}
