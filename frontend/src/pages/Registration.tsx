import { useState } from "react";
import type { FormEvent } from "react";

import {
  registerTeam,
} from "../api/registrationApi";
import PublicFooter from "../components/PublicFooter";
import PublicHeader from "../components/PublicHeader";

const EVENT_ID =
  "00000000-0000-0000-0000-000000000002";

function Registration() {
  const [player1Name, setPlayer1Name] = useState("");
  const [player1Phone, setPlayer1Phone] = useState("");

  const [player2Name, setPlayer2Name] = useState("");
  const [player2Phone, setPlayer2Phone] = useState("");

  const [teamName, setTeamName] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [registrationId, setRegistrationId] =
    useState("");

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault();

    setError("");
    setRegistrationId("");
    setLoading(true);

    try {
      const result = await registerTeam(EVENT_ID, {
        player1: {
          name: player1Name,
          phone: player1Phone,
        },
        player2: {
          name: player2Name,
          phone: player2Phone,
        },
        teamName,
      });

      setRegistrationId(result.registrationId);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Something went wrong",
      );
    } finally {
      setLoading(false);
    }
  }

  if (registrationId) {
    return (
      <div className="min-h-screen bg-slate-950 text-white">
        <PublicHeader />
        <main className="px-5 py-10 sm:px-8">
          <div className="mx-auto max-w-lg rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-2xl shadow-slate-950/40 sm:p-8">

          <div className="text-center">
            <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full border border-emerald-400/30 bg-emerald-500/10 text-2xl">✓</div>

            <h1 className="mt-4 text-2xl font-bold text-white">
              Registration Submitted!
            </h1>

            <p className="mt-2 text-slate-300">
              Your Men's Doubles team has been
              registered successfully.
            </p>
          </div>

          <div className="mt-6 rounded-xl border border-slate-700 bg-slate-950 p-4">
            <p className="text-sm text-slate-400">
              Registration ID
            </p>

            <p className="mt-1 break-all font-mono text-sm text-slate-200">
              {registrationId}
            </p>
          </div>

          <div className="mt-4 text-center">
            <span className="rounded-full border border-amber-400/30 bg-amber-500/10 px-4 py-2 text-sm text-amber-200">
              Pending confirmation
            </span>
          </div>

          <p className="mt-6 text-center text-sm text-slate-400">
            Your registration is pending admin confirmation. Check the{" "}
            <a
              className="font-semibold text-amber-300 underline underline-offset-2 hover:text-amber-200"
              href="/teams"
            >
              View Teams
            </a>{" "}
            page to confirm your registration status.
          </p>

          <a
            className="mt-6 block w-full rounded-lg bg-white px-4 py-3 text-center font-semibold text-slate-950 transition hover:bg-slate-200"
            href="/"
          >
            Return to homepage
          </a>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <PublicHeader />

      <main className="px-5 py-10 sm:px-8">
      <div className="mx-auto max-w-2xl">
        <div className="mb-8">
          <p className="text-sm font-medium uppercase tracking-[0.2em] text-amber-300">
            Badminton Tournament 2026
          </p>
          <h1 className="mt-3 text-3xl font-black tracking-tight sm:text-4xl">
            Register your team
          </h1>
          <p className="mt-3 max-w-xl text-slate-300">
            Share your player details to reserve a place in the Men's Doubles tournament.
          </p>
        </div>

        <form
          onSubmit={handleSubmit}
          className="space-y-6 rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-2xl shadow-slate-950/40 sm:p-8"
        >

          {/* Player 1 */}
          <div>
            <h2 className="text-lg font-semibold text-white">
              Player 1
            </h2>

            <div className="mt-3 space-y-3">

              <input
                required
                value={player1Name}
                onChange={(e) =>
                  setPlayer1Name(e.target.value)
                }
                placeholder="Full name"
                aria-label="Player 1 full name"
                className="w-full rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
              />

              <input
                required
                type="tel"
                inputMode="numeric"
                pattern="[6-9][0-9]{9}"
                value={player1Phone}
                onChange={(e) =>
                  setPlayer1Phone(e.target.value)
                }
                placeholder="Phone number"
                aria-label="Player 1 phone number"
                className="w-full rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
              />

            </div>
          </div>

          {/* Player 2 */}
          <div>
            <h2 className="text-lg font-semibold text-white">
              Player 2
            </h2>

            <div className="mt-3 space-y-3">

              <input
                required
                value={player2Name}
                onChange={(e) =>
                  setPlayer2Name(e.target.value)
                }
                placeholder="Full name"
                aria-label="Player 2 full name"
                className="w-full rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
              />

              <input
                type="tel"
                inputMode="numeric"
                pattern="[6-9][0-9]{9}"
                value={player2Phone}
                onChange={(e) =>
                  setPlayer2Phone(e.target.value)
                }
                placeholder="Phone number (optional)"
                aria-label="Player 2 phone number, optional"
                className="w-full rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
              />

            </div>
          </div>

          {/* Team name */}
          <div>
            <label className="text-sm font-medium text-white">
              Team Name
              <span className="ml-1 text-slate-400">
                (optional)
              </span>
            </label>

            <input
              value={teamName}
              onChange={(e) =>
                setTeamName(e.target.value)
              }
              placeholder="e.g. Smash Brothers"
                className="mt-2 w-full rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-white outline-none placeholder:text-slate-500 focus:border-amber-400 focus:ring-1 focus:ring-amber-400"
            />
          </div>

          {error && (
            <div className="rounded-lg border border-red-400/30 bg-red-950/40 p-3 text-sm text-red-200">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-white px-4 py-3 font-semibold text-slate-950 transition hover:bg-slate-200 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading
              ? "Registering..."
              : "Register Team"}
          </button>

        </form>
      </div>
      </main>

      <PublicFooter />
    </div>
  );
}

export default Registration;