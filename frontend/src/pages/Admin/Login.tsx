import { useState } from "react";
import type { FormEvent } from "react";
import { useNavigate } from "react-router-dom";

import { login } from "../../api/authApi";
import PublicFooter from "../../components/PublicFooter";
import PublicHeader from "../../components/PublicHeader";

export default function Login() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login(email, password);
      navigate("/admin/registrations", { replace: true });
    } catch (loginError) {
      setError(
        loginError instanceof Error
          ? loginError.message
          : "Unable to sign in",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen flex-col bg-slate-950 text-white">
      <PublicHeader />
      <main className="flex flex-1 items-center justify-center px-4 py-10">
      <form
        className="w-full max-w-md rounded-2xl border border-slate-800 bg-slate-900 p-8 shadow-2xl shadow-slate-950/40"
        onSubmit={handleSubmit}
      >
        <p className="text-sm font-medium uppercase tracking-[0.2em] text-amber-300">
          ShuttleHub administration
        </p>
        <h1 className="mt-3 text-3xl font-bold text-white">Admin sign in</h1>

        <input
          className="mt-6 w-full rounded-lg border border-slate-700 bg-slate-950 p-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
          onChange={(event) => setEmail(event.target.value)}
          placeholder="Email"
          required
          type="email"
          value={email}
        />

        <input
          className="mt-4 w-full rounded-lg border border-slate-700 bg-slate-950 p-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
          onChange={(event) => setPassword(event.target.value)}
          placeholder="Password"
          required
          type="password"
          value={password}
        />

        {error && (
          <p className="mt-4 text-sm text-red-300">{error}</p>
        )}

        <button
          className="mt-6 w-full rounded-lg bg-white p-3 font-semibold text-slate-950 transition hover:bg-slate-200 disabled:cursor-not-allowed disabled:opacity-50"
          disabled={loading}
          type="submit"
        >
          {loading ? "Signing in..." : "Sign in"}
        </button>
      </form>
      </main>
      <PublicFooter />
    </div>
  );
}