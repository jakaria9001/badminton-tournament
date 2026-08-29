import { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getEvent } from "../../api/eventApi";
import {
    getAdminProfile,
    getAdminToken,
    logout,
    type AdminProfile,
} from "../../api/authApi";
import PublicFooter from "../../components/PublicFooter";
import PublicHeader from "../../components/PublicHeader";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL;

interface Registration {
  id: string;
  teamId: string;
  teamName: string;
  player1Name: string;
  player1Phone: string;
  player2Name: string;
  player2Phone: string;
  status: string;
  registeredAt: string;
}

function Registrations() {

    const navigate = useNavigate();
    const { eventId: routeEventId } = useParams();

    const handleUnauthorized = useCallback(() => {
        logout();
        navigate("/admin/login", { replace: true });
    }, [navigate]);

    const [registrations, setRegistrations] =
        useState<Registration[]>([]);

    const [loading, setLoading] =
        useState(true);

    const [profile, setProfile] =
        useState<AdminProfile | null>(null);

    const [registrationStatus, setRegistrationStatus] =
        useState("REGISTRATION_CLOSED");
    const eventId = routeEventId ?? profile?.eventId ?? "";

    const sortedRegistrations = [...registrations].sort((a, b) => {
        const priority: Record<string, number> = {
            PENDING: 0,
            CONFIRMED: 1,
            REJECTED: 2,
        };

        return (priority[a.status] ?? 99) - (priority[b.status] ?? 99);
    });

    const summary = {
        total: registrations.length,
        pending: registrations.filter((registration) => registration.status === "PENDING").length,
        confirmed: registrations.filter((registration) => registration.status === "CONFIRMED").length,
        rejected: registrations.filter((registration) => registration.status === "REJECTED").length,
    };

    const loadRegistrations = useCallback(async () => {
        if (!eventId) {
            setLoading(false);
            return;
        }

        try {
            const response = await fetch(
                `${API_BASE_URL}/api/v1/admin/events/${eventId}/registrations`,
                {
                    headers: {
                    Authorization: `Bearer ${getAdminToken()}`,
                    },
                },
            );

            if (response.status === 401 || response.status === 403) {
                    handleUnauthorized();
                    return;
            }

            if (!response.ok) {
                throw new Error(
                "Failed to load registrations",
                );
            }

            const data = await response.json();

            setRegistrations(data ?? []);
        } finally {
            setLoading(false);
        }
    }, [eventId, handleUnauthorized]);

    useEffect(() => {
        async function loadProfile() {
            try {
            const adminProfile = await getAdminProfile();
            setProfile(adminProfile);
            } catch {
            handleUnauthorized();
            }
        }

        void loadProfile();
    }, [handleUnauthorized, navigate]);

    useEffect(() => {
        void loadRegistrations();
    }, [eventId, loadRegistrations]);

    useEffect(() => {
        if (!eventId) {
            return;
        }

        void getEvent(eventId).then((event) => {
            setRegistrationStatus(event.status);
        }).catch(() => undefined);
    }, [eventId]);

    async function toggleRegistrationStatus() {
        const nextStatus = registrationStatus === "REGISTRATION_OPEN"
            ? "REGISTRATION_CLOSED"
            : "REGISTRATION_OPEN";
        if (!eventId) {
            return;
        }

        const response = await fetch(
            `${API_BASE_URL}/api/v1/admin/events/${eventId}/registration-status`,
            {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${getAdminToken()}`,
                },
                body: JSON.stringify({ status: nextStatus }),
            },
        );

        if (response.status === 401 || response.status === 403) {
            handleUnauthorized();
            return;
        }
        if (!response.ok) {
            alert(await response.text());
            return;
        }
        setRegistrationStatus(nextStatus);
    }

    async function updateStatus(
        registrationId: string,
        status: string,
    ) {
        const response = await fetch(
        `${API_BASE_URL}/api/v1/admin/registrations/${registrationId}/status`,
        {
            method: "PUT",
            headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${getAdminToken()}`,
            },
            body: JSON.stringify({
            status,
            }),
        },
        );

        if (response.status === 401 || response.status === 403) {
            handleUnauthorized();
            return;
        }
        
        if (!response.ok) {
            alert("Failed to update registration");
            return;
        }

        await loadRegistrations();
    }

    async function withdrawTeam(registrationId: string) {
        const confirmed = window.confirm(
            "Withdraw this team registration? It will leave the public teams list, but its history will be retained.",
        );

        if (!confirmed) {
            return;
        }

        const response = await fetch(
            `${API_BASE_URL}/api/v1/admin/registrations/${registrationId}/withdraw`,
            {
                method: "PUT",
                headers: {
                    Authorization: `Bearer ${getAdminToken()}`,
                },
            },
        );

        if (response.status === 401 || response.status === 403) {
            handleUnauthorized();
            return;
        }

        if (!response.ok) {
            alert("Failed to withdraw team");
            return;
        }

        await loadRegistrations();
    }

    if (loading) {
        return (
        <div className="flex min-h-screen flex-col bg-slate-100">
            <PublicHeader />
            <main className="flex flex-1 items-center justify-center p-10 text-slate-700">
                Loading registrations...
            </main>
            <PublicFooter />
        </div>
        );
    }

    return (
        <div className="flex min-h-screen flex-col bg-slate-100">
            <PublicHeader />
            <main className="flex-1 px-4 py-8">
            <div className="mx-auto max-w-5xl">
                <div className="mb-8 rounded-2xl bg-slate-900 p-6 text-white shadow-lg">
                    <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
                        <div>
                            <p className="text-sm uppercase tracking-[0.18em] text-slate-300">
                                Admin
                            </p>
                            <h1 className="mt-2 text-3xl font-bold">
                                Men's Doubles Registrations
                            </h1>
                        </div>

                        <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
                            <div className="hidden text-right sm:block">
                                <p className="text-sm font-semibold text-white">
                                    {profile?.name ?? "Loading..."}
                                </p>
                                <p className="text-xs text-slate-300">
                                    {profile?.role ?? ""}
                                </p>
                            </div>

                            <div className="flex items-center gap-2">
                                <button
                                    className={`rounded-lg px-4 py-2 text-sm font-semibold transition ${registrationStatus === "REGISTRATION_OPEN" ? "border border-amber-400/40 bg-amber-500/10 text-amber-200 hover:bg-amber-500/20" : "border border-emerald-400/40 bg-emerald-500/10 text-emerald-200 hover:bg-emerald-500/20"}`}
                                    onClick={() => void toggleRegistrationStatus()}
                                    type="button"
                                >
                                    {registrationStatus === "REGISTRATION_OPEN" ? "Stop registrations" : "Open registrations"}
                                </button>
                                <button
                                    className="rounded-lg border border-white/20 bg-white/5 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/10"
                                    onClick={() => navigate("/admin/draw")}
                                    type="button"
                                >
                                    Draw
                                </button>

                                <button
                                    className="rounded-lg border border-white/20 bg-white/5 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/10"
                                    onClick={() => void loadRegistrations()}
                                    type="button"
                                >
                                    Refresh
                                </button>

                                <button
                                    className="rounded-lg border border-red-400/40 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-200 transition hover:bg-red-500 hover:text-white"
                                    onClick={handleUnauthorized}
                                    type="button"
                                >
                                    Sign out
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="mb-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <StatCard label="Total" value={summary.total} tone="slate" />
                    <StatCard label="Pending" value={summary.pending} tone="yellow" />
                    <StatCard label="Confirmed" value={summary.confirmed} tone="green" />
                    <StatCard label="Rejected" value={summary.rejected} tone="red" />
                </div>

                {sortedRegistrations.length === 0 ? (
                    <div className="rounded-2xl bg-white p-8 text-center shadow-sm">
                        <p className="text-lg font-semibold text-slate-800">
                            No registrations yet
                        </p>
                        <p className="mt-2 text-sm text-slate-500">
                            New team submissions will appear here.
                        </p>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {sortedRegistrations.map((registration) => {
                            const statusStyle: Record<string, string> = {
                                PENDING: "bg-yellow-100 text-yellow-800",
                                CONFIRMED: "bg-green-100 text-green-800",
                                REJECTED: "bg-red-100 text-red-800",
                            };

                            return (
                                <div
                                    key={registration.id}
                                    className="rounded-2xl bg-white p-6 shadow-sm"
                                >
                                    <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
                                        <div>
                                            <h2 className="text-lg font-bold">
                                                {registration.teamName}
                                            </h2>

                                            <div className="mt-2 space-y-1 text-sm text-slate-600">
                                                <p>
                                                    {registration.player1Name}
                                                    {" · "}
                                                    {registration.player1Phone}
                                                </p>

                                                <p>
                                                    {registration.player2Name}
                                                    {" · "}
                                                    {registration.player2Phone}
                                                </p>
                                            </div>
                                        </div>

                                        <div className="flex flex-wrap items-center gap-3">
                                            <span
                                                className={`rounded-full px-3 py-1 text-xs font-medium ${statusStyle[registration.status] ?? "bg-slate-100 text-slate-800"}`}
                                            >
                                                {registration.status}
                                            </span>

                                            {registration.status === "PENDING" && (
                                                <>
                                                    <button
                                                        onClick={() =>
                                                            updateStatus(
                                                                registration.id,
                                                                "CONFIRMED",
                                                            )
                                                        }
                                                        className="rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700"
                                                        type="button"
                                                    >
                                                        Confirm
                                                    </button>

                                                    <button
                                                        onClick={() =>
                                                            updateStatus(
                                                                registration.id,
                                                                "REJECTED",
                                                            )
                                                        }
                                                        className="rounded-lg bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700"
                                                        type="button"
                                                    >
                                                        Reject
                                                    </button>
                                                </>
                                            )}

                                            <button
                                                onClick={() => withdrawTeam(registration.id)}
                                                className="rounded-lg border border-red-200 bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 hover:bg-red-100"
                                                type="button"
                                            >
                                                Withdraw team
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
            </main>
            <PublicFooter />
        </div>
    );
}

interface StatCardProps {
    label: string;
    value: number;
    tone: "slate" | "yellow" | "green" | "red";
}

function StatCard({ label, value, tone }: StatCardProps) {
    const toneStyles: Record<StatCardProps["tone"], string> = {
        slate: "bg-slate-900 text-white",
        yellow: "bg-yellow-100 text-yellow-800",
        green: "bg-green-100 text-green-800",
        red: "bg-red-100 text-red-800",
    };

    return (
        <div className={`rounded-2xl p-4 shadow-sm ${toneStyles[tone]}`}>
            <p className="text-sm opacity-80">{label}</p>
            <p className="mt-2 text-3xl font-bold">{value}</p>
        </div>
    );
}

export default Registrations;