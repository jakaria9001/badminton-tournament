import { useEffect, useState, type ReactNode } from "react";
import { Link } from "react-router-dom";
import { getAdminProfile, type AdminProfile } from "../api/authApi";

interface PublicHeaderProps {
  tournamentName?: string;
  adminActions?: ReactNode;
}

export default function PublicHeader({ tournamentName, adminActions }: PublicHeaderProps) {
  const [profile, setProfile] = useState<AdminProfile | null>(null);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    let active = true;

    void getAdminProfile()
      .then((currentProfile) => {
        if (active) {
          setProfile(currentProfile);
        }
      })
      .catch(() => {
        if (active) {
          setProfile(null);
        }
      });

    return () => {
      active = false;
    };
  }, []);

  const dashboardPath = profile?.role === "SUPER_ADMIN" ? "/admin/superadmin" : "/admin";

  return (
    <header className="border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-sm">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-5 py-4 sm:px-8">
        <Link to="/" className="flex items-center gap-3" aria-label="ShuttleHub home">
          <span className="flex shrink-0 rounded-lg border border-white bg-white p-1 shadow-sm">
            <img
              src="/shuttlehub_logo.svg"
              alt="ShuttleHub logo"
              className="block h-9 w-auto"
            />
          </span>
          <div>
            <div className="text-lg font-bold tracking-tight text-white">ShuttleHub</div>
            <div className="text-[10px] uppercase tracking-[0.22em] text-slate-400">Tournaments</div>
          </div>
        </Link>

        <button
          aria-controls="mobile-navigation"
          aria-expanded={isMenuOpen}
          aria-label="Open navigation menu"
          className="flex h-10 w-10 items-center justify-center rounded-lg border border-slate-700 text-slate-200 transition hover:border-amber-400 hover:text-white"
          onClick={() => setIsMenuOpen((open) => !open)}
          type="button"
        >
          <span className="flex w-5 flex-col gap-1.5" aria-hidden="true">
            <span className="h-0.5 w-full bg-current" />
            <span className="h-0.5 w-full bg-current" />
            <span className="h-0.5 w-full bg-current" />
          </span>
        </button>
      </div>

      {isMenuOpen && (
        <nav
          className="border-t border-slate-800 bg-slate-950 px-5 py-4 text-sm text-slate-200"
          id="mobile-navigation"
        >
          {profile && (
            <div className="border-b border-slate-800 pb-4">
              {tournamentName && <p className="text-xs font-semibold uppercase tracking-[0.14em] text-amber-300">{tournamentName}</p>}
              <p className="mt-1 font-semibold text-white">{profile.name}</p>
              <p className="mt-1 text-xs text-slate-400">{profile.role.replaceAll("_", " ")}</p>
            </div>
          )}
          <div className="flex flex-col py-2">
            <Link className="rounded-lg px-3 py-2.5 transition hover:bg-slate-900 hover:text-white" onClick={() => setIsMenuOpen(false)} to="/">Events</Link>
            <Link className="rounded-lg px-3 py-2.5 transition hover:bg-slate-900 hover:text-white" onClick={() => setIsMenuOpen(false)} to={profile ? dashboardPath : "/admin/login"}>
              {profile ? "Dashboard" : "Login"}
            </Link>
          </div>
          {profile && adminActions && (
            <div className="flex flex-col gap-2 border-t border-slate-800 pt-4">
              {adminActions}
            </div>
          )}
        </nav>
      )}
    </header>
  );
}
