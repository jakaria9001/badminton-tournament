import { getAdminToken } from "../api/authApi";

export default function PublicHeader() {
  const hasAdminToken = Boolean(getAdminToken());

  return (
    <header className="border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-sm">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-5 py-4 sm:px-8">
        <a href="/" className="flex items-center gap-3" aria-label="ShuttleHub home">
          <span className="flex shrink-0 rounded-lg border border-white bg-white p-1 shadow-sm">
            <img
              src="/shuttlehub_logo.svg"
              alt="ShuttleHub logo"
              className="block h-9 w-auto"
            />
          </span>
          <div>
            <div className="text-lg font-bold tracking-tight text-white">ShuttleHub</div>
            <div className="text-[10px] uppercase tracking-[0.22em] text-slate-400">Men's Doubles</div>
          </div>
        </a>

        <nav className="hidden items-center gap-6 text-sm text-slate-300 md:flex">
          <a href="/" className="transition hover:text-white">Home</a>
          <a href="/teams" className="transition hover:text-white">Teams</a>
          <a href="/register" className="transition hover:text-white">Register</a>
          <a href={hasAdminToken ? "/admin/registrations" : "/admin/login"} className="transition hover:text-white">
            {hasAdminToken ? "Dashboard" : "Login"}
          </a>
        </nav>
      </div>
    </header>
  );
}
