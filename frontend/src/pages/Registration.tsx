import { useState } from "react";
import type { FormEvent } from "react";

import {
  registerTeam,
} from "../api/registrationApi";

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
      <div className="min-h-screen bg-slate-100 px-4 py-10">
        <div className="mx-auto max-w-lg rounded-2xl bg-white p-8 shadow-lg">

          <div className="text-center">
            <div className="text-5xl">
              🎉
            </div>

            <h1 className="mt-4 text-2xl font-bold">
              Registration Submitted!
            </h1>

            <p className="mt-2 text-slate-600">
              Your Men's Doubles team has been
              registered successfully.
            </p>
          </div>

          <div className="mt-6 rounded-xl bg-slate-100 p-4">
            <p className="text-sm text-slate-500">
              Registration ID
            </p>

            <p className="mt-1 break-all font-mono text-sm">
              {registrationId}
            </p>
          </div>

          <div className="mt-4 text-center">
            <span className="rounded-full bg-yellow-100 px-4 py-2 text-sm text-yellow-800">
              Pending confirmation
            </span>
          </div>

          <p className="mt-6 text-center text-sm text-slate-600">
            Your registration is pending admin confirmation. Check the{" "}
            <a
              className="font-semibold text-slate-900 underline underline-offset-2 hover:text-slate-600"
              href="/teams"
            >
              View Teams
            </a>{" "}
            page to confirm your registration status.
          </p>

          <a
            className="mt-6 block w-full rounded-lg bg-slate-900 px-4 py-3 text-center font-semibold text-white transition hover:bg-slate-800"
            href="/"
          >
            Return to homepage
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-100 px-4 py-10">

      <div className="mx-auto max-w-lg">

        <div className="mb-8 text-center">
          <div className="text-4xl">
            🏸
          </div>

          <h1 className="mt-3 text-3xl font-bold text-slate-900">
            Badminton Open 2026
          </h1>

          <p className="mt-2 text-slate-600">
            Men's Doubles Registration
          </p>
        </div>

        <form
          onSubmit={handleSubmit}
          className="space-y-6 rounded-2xl bg-white p-6 shadow-lg"
        >

          {/* Player 1 */}
          <div>
            <h2 className="text-lg font-semibold">
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
                className="w-full rounded-lg border border-slate-300 px-4 py-3 outline-none focus:border-slate-900"
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
                className="w-full rounded-lg border border-slate-300 px-4 py-3 outline-none focus:border-slate-900"
              />

            </div>
          </div>

          {/* Player 2 */}
          <div>
            <h2 className="text-lg font-semibold">
              Player 2
              <span className="ml-2 text-sm font-normal text-slate-400">
                (phone optional)
              </span>
            </h2>

            <div className="mt-3 space-y-3">

              <input
                required
                value={player2Name}
                onChange={(e) =>
                  setPlayer2Name(e.target.value)
                }
                placeholder="Full name"
                className="w-full rounded-lg border border-slate-300 px-4 py-3 outline-none focus:border-slate-900"
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
                className="w-full rounded-lg border border-slate-300 px-4 py-3 outline-none focus:border-slate-900"
              />

            </div>
          </div>

          {/* Team name */}
          <div>
            <label className="text-sm font-medium">
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
              className="mt-2 w-full rounded-lg border border-slate-300 px-4 py-3 outline-none focus:border-slate-900"
            />
          </div>

          {error && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-slate-900 px-4 py-3 font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading
              ? "Registering..."
              : "Register Team"}
          </button>

        </form>
      </div>
    </div>
  );
}

export default Registration;