import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { listAdminEvents } from "../api/eventApi";
import type { EventInfo } from "../types/event";
import { getAdminProfile, getAdminToken, logout } from "../api/authApi";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: string;
  eventId?: string;
}

export default function SuperAdminDashboard() {
  const navigate = useNavigate();
  const [events, setEvents] = useState<EventInfo[]>([]);
  const [admins, setAdmins] = useState<AdminUser[]>([]);
  const [name, setName] = useState("");
  const [venueName, setVenueName] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [maxTeams, setMaxTeams] = useState("");
  const [status, setStatus] = useState("DRAFT");
  const [assignedAdminId, setAssignedAdminId] = useState("");
  const [adminName, setAdminName] = useState("");
  const [adminEmail, setAdminEmail] = useState("");
  const [adminPassword, setAdminPassword] = useState("");
  const [adminRole, setAdminRole] = useState("ADMIN");
  const [adminEventId, setAdminEventId] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState<string | null>(null);

  const loadDashboard = useCallback(async () => {
    const [eventData, adminResponse] = await Promise.all([
      listAdminEvents(),
      fetch(`${API_BASE_URL}/api/v1/admin/superadmin/admins`, {
        headers: { Authorization: `Bearer ${getAdminToken()}` },
      }),
    ]);

    if (!adminResponse.ok) {
      throw new Error("Unable to load admin roster");
    }

    setEvents(eventData);
    setAdmins((await adminResponse.json()) ?? []);
  }, []);

  useEffect(() => {
    async function bootstrap() {
      try {
        const profile = await getAdminProfile();
        if (profile.role !== "SUPER_ADMIN") {
          navigate("/admin", { replace: true });
          return;
        }
        await loadDashboard();
      } catch (loadError) {
        if (loadError instanceof Error && loadError.message !== "Session expired") {
          setError(loadError.message);
        }
        logout();
        navigate("/admin/login", { replace: true });
      } finally {
        setLoading(false);
      }
    }

    void bootstrap();
  }, [loadDashboard, navigate]);

  async function createEvent() {
    setBusy("event");
    setError("");
    setSuccess("");
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/admin/superadmin/events`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAdminToken()}` },
        body: JSON.stringify({
          name,
          venueName,
          startDate,
          endDate,
          maxTeams: maxTeams ? Number(maxTeams) : null,
          status,
          assignedAdminId: assignedAdminId || undefined,
        }),
      });
      if (!response.ok) throw new Error(await response.text());
      setName("");
      setVenueName("");
      setStartDate("");
      setEndDate("");
      setMaxTeams("");
      setStatus("DRAFT");
      setAssignedAdminId("");
      await loadDashboard();
      setSuccess("Tournament created.");
    } catch (createError) {
      setError(createError instanceof Error ? createError.message : "Unable to create event");
    } finally {
      setBusy(null);
    }
  }

  async function createAdminAccount() {
    setBusy("admin");
    setError("");
    setSuccess("");
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/admin/superadmin/admins`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAdminToken()}` },
        body: JSON.stringify({
          name: adminName,
          email: adminEmail,
          password: adminPassword,
          role: adminRole,
          eventId: adminRole === "ADMIN" ? adminEventId : undefined,
        }),
      });
      if (!response.ok) throw new Error(await response.text());
      setAdminName("");
      setAdminEmail("");
      setAdminPassword("");
      setAdminRole("ADMIN");
      setAdminEventId("");
      await loadDashboard();
      setSuccess("Admin account created.");
    } catch (createError) {
      setError(createError instanceof Error ? createError.message : "Unable to create admin");
    } finally {
      setBusy(null);
    }
  }

  async function deleteEvent(event: EventInfo) {
    if (!window.confirm(`Delete ${event.name}? This removes its tournament data.`)) return;
    setError("");
    setSuccess("");
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/admin/superadmin/events/${event.id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${getAdminToken()}` },
      });
      if (!response.ok) throw new Error(await response.text());
      await loadDashboard();
      setSuccess("Tournament deleted.");
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Unable to delete event");
    }
  }

  function signOut() {
    logout();
    navigate("/admin/login", { replace: true });
  }

  if (loading) {
    return (
      <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900">
        <PublicHeader />
        <main className="flex flex-1 items-center justify-center p-10 text-slate-700">Loading platform dashboard...</main>
        <PublicFooter />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col bg-slate-100 text-slate-900">
      <PublicHeader />
      <main className="flex-1 px-4 py-8">
        <div className="mx-auto max-w-6xl">
          <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
            <div>
              <p className="text-sm font-bold uppercase tracking-[0.18em] text-slate-500">ShuttleHub · SuperAdmin</p>
              <h1 className="mt-2 text-3xl font-black text-slate-950">Platform Control</h1>
              <p className="mt-2 text-sm text-slate-500">Manage tournaments and admin roles across the platform.</p>
            </div>
            <button className="rounded-lg border border-red-200 bg-red-50 px-4 py-2 text-sm font-bold text-red-700" onClick={signOut} type="button">Sign out</button>
          </div>

          {error && <p className="mb-5 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900" role="alert">{error}</p>}
          {success && <p className="mb-5 rounded-lg border border-emerald-200 bg-emerald-50 p-4 text-sm text-emerald-900" role="status">{success}</p>}

          <div className="grid gap-4 sm:grid-cols-3">
            <Metric label="Total events" value={events.length} />
            <Metric label="Open registrations" value={events.filter((event) => event.status === "REGISTRATION_OPEN").length} />
            <Metric label="Registered teams" value={events.reduce((total, event) => total + event.registeredTeams, 0)} />
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-[1.1fr_1.4fr]">
            <section className="rounded-2xl bg-white p-6 shadow-sm">
              <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">New event</p>
              <h2 className="mt-2 text-xl font-black text-slate-950">Create tournament</h2>
              <div className="mt-5 space-y-3">
                <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setName(event.target.value)} placeholder="Event name" value={name} />
                <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setVenueName(event.target.value)} placeholder="Venue" value={venueName} />
                <div className="grid gap-3 sm:grid-cols-2">
                  <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setStartDate(event.target.value)} required type="date" value={startDate} />
                  <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setEndDate(event.target.value)} required type="date" value={endDate} />
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setMaxTeams(event.target.value)} placeholder="Max teams" type="number" value={maxTeams} />
                  <select className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setStatus(event.target.value)} value={status}>
                    <option value="DRAFT">Draft</option>
                    <option value="REGISTRATION_OPEN">Registration open</option>
                    <option value="REGISTRATION_CLOSED">Registration closed</option>
                  </select>
                </div>
                <label className="block text-sm font-semibold text-slate-700">
                  Assign admin to this event
                  <select className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAssignedAdminId(event.target.value)} value={assignedAdminId}>
                    <option value="">No admin assignment</option>
                    {admins.filter((admin) => admin.role === "ADMIN").map((admin) => (
                      <option key={admin.id} value={admin.id}>{admin.name} ({admin.email})</option>
                    ))}
                  </select>
                </label>
                <button className="w-full rounded-lg bg-slate-950 px-4 py-2.5 text-sm font-bold text-white disabled:cursor-not-allowed disabled:opacity-60" disabled={busy !== null || !name || !venueName || !startDate || !endDate} onClick={() => void createEvent()} type="button">
                  {busy === "event" ? "Creating..." : "Create tournament"}
                </button>
              </div>
            </section>

            <section className="rounded-2xl bg-white p-6 shadow-sm">
              <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">New admin</p>
              <h2 className="mt-2 text-xl font-black text-slate-950">Create administrator</h2>
              <div className="mt-5 space-y-3">
                <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAdminName(event.target.value)} placeholder="Full name" value={adminName} />
                <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAdminEmail(event.target.value)} placeholder="Email" type="email" value={adminEmail} />
                <input className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAdminPassword(event.target.value)} placeholder="Password" type="password" value={adminPassword} />
                <div className="grid gap-3 sm:grid-cols-2">
                  <select className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAdminRole(event.target.value)} value={adminRole}>
                    <option value="ADMIN">Admin</option>
                    <option value="SUPER_ADMIN">Super admin</option>
                  </select>
                  {adminRole === "ADMIN" && (
                    <select className="w-full rounded-lg border border-slate-300 px-3 py-2" onChange={(event) => setAdminEventId(event.target.value)} value={adminEventId}>
                      <option value="">Select assigned event</option>
                      {events.map((event) => (
                        <option key={event.id} value={event.id}>{event.name}</option>
                      ))}
                    </select>
                  )}
                </div>
                <button className="w-full rounded-lg bg-amber-500 px-4 py-2.5 text-sm font-bold text-slate-950 disabled:cursor-not-allowed disabled:opacity-60" disabled={busy !== null || !adminName || !adminEmail || !adminPassword || (adminRole === "ADMIN" && !adminEventId)} onClick={() => void createAdminAccount()} type="button">
                  {busy === "admin" ? "Creating..." : "Create admin"}
                </button>
              </div>
            </section>
          </div>

          <section className="mt-8 rounded-2xl bg-white p-6 shadow-sm">
            <div className="mb-4 flex items-center justify-between gap-3">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Admin roster</p>
                <h2 className="mt-2 text-xl font-black text-slate-950">Assigned administrators</h2>
              </div>
            </div>
            <div className="grid gap-3 md:grid-cols-2">
              {admins.map((admin) => (
                <div className="rounded-xl border border-slate-200 bg-slate-50 p-4" key={admin.id}>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-bold text-slate-950">{admin.name}</p>
                      <p className="text-sm text-slate-600">{admin.email}</p>
                    </div>
                    <span className={`rounded-full px-2.5 py-1 text-[10px] font-bold uppercase tracking-wide ${admin.role === "SUPER_ADMIN" ? "bg-violet-100 text-violet-700" : "bg-emerald-100 text-emerald-700"}`}>{admin.role}</span>
                  </div>
                  <p className="mt-3 text-sm text-slate-600">{admin.role === "ADMIN" ? `Event: ${events.find((event) => event.id === admin.eventId)?.name ?? "Unassigned"}` : "Platform-wide access"}</p>
                </div>
              ))}
              {admins.length === 0 && <p className="text-sm text-slate-500">No admin accounts yet.</p>}
            </div>
          </section>

          <section className="mt-8 rounded-2xl bg-white p-6 shadow-sm">
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Tournaments</p>
            <h2 className="mt-2 text-xl font-black text-slate-950">All events</h2>
            <div className="mt-4 space-y-3">
              {events.map((event) => (
                <div className="flex flex-col gap-3 rounded-xl border border-slate-200 bg-slate-50 p-4 sm:flex-row sm:items-center sm:justify-between" key={event.id}>
                  <div>
                    <p className="font-bold text-slate-950">{event.name}</p>
                    <p className="text-sm text-slate-600">{event.status} · {event.registeredTeams} teams registered</p>
                  </div>
                  <button className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm font-semibold text-red-700" onClick={() => void deleteEvent(event)} type="button">Delete</button>
                </div>
              ))}
              {events.length === 0 && <p className="text-sm text-slate-500">No tournaments created yet.</p>}
            </div>
          </section>
        </div>
      </main>
      <PublicFooter />
    </div>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-2xl bg-white p-5 shadow-sm">
      <p className="text-3xl font-black text-slate-950">{value}</p>
      <p className="mt-2 text-xs font-bold uppercase tracking-wide text-slate-500">{label}</p>
    </div>
  );
}
