import { useState } from "react";
import type { FormEvent } from "react";
import { useNavigate } from "react-router-dom";

import { login } from "../../api/authApi";

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
    <main className="flex min-h-screen items-center justify-center bg-slate-950 px-4">
      <form
        className="w-full max-w-md rounded-2xl bg-white p-8 shadow-xl"
        onSubmit={handleSubmit}
      >
        <h1 className="text-3xl font-bold">Admin sign in</h1>

        <input
          className="mt-6 w-full rounded-lg border p-3"
          onChange={(event) => setEmail(event.target.value)}
          placeholder="Email"
          required
          type="email"
          value={email}
        />

        <input
          className="mt-4 w-full rounded-lg border p-3"
          onChange={(event) => setPassword(event.target.value)}
          placeholder="Password"
          required
          type="password"
          value={password}
        />

        {error && (
          <p className="mt-4 text-sm text-red-700">{error}</p>
        )}

        <button
          className="mt-6 w-full rounded-lg bg-slate-950 p-3 font-semibold text-white"
          disabled={loading}
          type="submit"
        >
          {loading ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </main>
  );
}