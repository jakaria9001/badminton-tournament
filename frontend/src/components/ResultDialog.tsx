import { useMemo, useState } from "react";

interface Team {
  id: string;
  teamName: string;
  player1: string;
  player2: string;
}

interface Match {
  id: string;
  roundId: string;
  team1: Team | null;
  team2: Team | null;
  status: string;
}

interface ResultDialogProps {
  match: Match;
  onClose: () => void;
  onSubmitted: () => Promise<void> | void;
  request: (url: string, options?: RequestInit) => Promise<Response>;
}

interface GameInput {
  gameNumber: number;
  winnerScore: number;
  loserScore: number;
}

const emptyGame = (gameNumber: number): GameInput => ({
  gameNumber,
  winnerScore: 0,
  loserScore: 0,
});

function ResultDialog({ match, onClose, onSubmitted, request }: ResultDialogProps) {
  const winnerOptions = useMemo(
    () => [match.team1, match.team2].filter((team): team is Team => Boolean(team)),
    [match.team1, match.team2],
  );

  const [winnerId, setWinnerId] = useState(match.team1?.id ?? winnerOptions[0]?.id ?? "");
  const [games, setGames] = useState<GameInput[]>([emptyGame(1)]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const winnerTeam = winnerOptions.find((team) => team.id === winnerId) ?? winnerOptions[0] ?? null;
  const loserTeam = winnerOptions.find((team) => team.id !== winnerId) ?? null;

  const updateGame = (gameNumber: number, field: "winnerScore" | "loserScore", value: string) => {
    const sanitized = value.replace(/[^\d]/g, "");
    const parsed = Number(sanitized);
    const clamped = Number.isFinite(parsed) ? Math.min(parsed, 30) : 0;

    setGames((current) =>
      current.map((game) =>
        game.gameNumber === gameNumber
          ? { ...game, [field]: clamped }
          : game,
      ),
    );
  };

  const addGame = () => {
    setGames((current) => [...current, emptyGame(current.length + 1)]);
  };

  const removeGame = (gameNumber: number) => {
    setGames((current) => {
      const filtered = current.filter((game) => game.gameNumber !== gameNumber);
      return filtered.length === 0 ? [emptyGame(1)] : filtered.map((game, index) => ({ ...game, gameNumber: index + 1 }));
    });
  };

  async function submit() {
    if (!match.team1 || !match.team2) {
      setError("Both teams must be available before submitting a result.");
      return;
    }

    if (!winnerId) {
      setError("Select the winning team.");
      return;
    }

    const normalizedGames = games.map((game) => {
      if (!winnerTeam || !loserTeam) {
        return {
          gameNumber: game.gameNumber,
          team1Score: 0,
          team2Score: 0,
        };
      }

      const winnerIsTeam1 = winnerTeam.id === match.team1?.id;

      return {
        gameNumber: game.gameNumber,
        team1Score: winnerIsTeam1 ? game.winnerScore : game.loserScore,
        team2Score: winnerIsTeam1 ? game.loserScore : game.winnerScore,
      };
    });

    try {
      setSubmitting(true);
      setError("");

      await request(`/api/v1/admin/matches/${match.id}/result`, {
        method: "POST",
        body: JSON.stringify({
          winnerTeamId: winnerId,
          games: normalizedGames,
        }),
      });

      await onSubmitted();
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "Unable to submit result");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4">
      <div className="w-full max-w-2xl rounded-2xl bg-white p-6 shadow-2xl">
        <div className="mb-5 flex items-center justify-between gap-3">
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">Match result</p>
            <h2 className="mt-1 text-2xl font-bold text-slate-900">
              {match.team1?.player1 ?? "TBD"} / {match.team1?.player2 ?? "TBD"} vs {match.team2?.player1 ?? "TBD"} / {match.team2?.player2 ?? "TBD"}
            </h2>
          </div>
          <button
            aria-label="Close result dialog"
            className="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-300 text-2xl leading-none text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
            onClick={onClose}
            type="button"
          >
            &times;
          </button>
        </div>

        <div className="mb-5">
          <label className="mb-2 block text-sm font-semibold text-slate-700">Winning team</label>
          <select
            className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2"
            value={winnerId}
            onChange={(event) => setWinnerId(event.target.value)}
          >
            {winnerOptions.map((team) => (
              <option key={team.id} value={team.id}>
                {team.player1} / {team.player2}
              </option>
            ))}
          </select>
        </div>

        <div className="mb-5 space-y-3">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold uppercase tracking-wide text-slate-600">Game scores</h3>
            <button
              className="rounded-md bg-slate-800 px-3 py-1.5 text-xs font-semibold text-white"
              onClick={addGame}
              type="button"
            >
              Add game
            </button>
          </div>

          {games.map((game) => (
            <div key={game.gameNumber} className="grid grid-cols-[auto_1fr_1fr_auto] items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 p-3">
              <span className="text-sm font-semibold text-slate-700">G{game.gameNumber}</span>
              <input
                aria-label={`Game ${game.gameNumber} winner score`}
                className="rounded-lg border border-slate-300 px-2 py-2 text-sm"
                inputMode="numeric"
                max={30}
                min={0}
                onChange={(event) => updateGame(game.gameNumber, "winnerScore", event.target.value)}
                type="text"
                value={game.winnerScore}
              />
              <input
                aria-label={`Game ${game.gameNumber} loser score`}
                className="rounded-lg border border-slate-300 px-2 py-2 text-sm"
                inputMode="numeric"
                max={30}
                min={0}
                onChange={(event) => updateGame(game.gameNumber, "loserScore", event.target.value)}
                type="text"
                value={game.loserScore}
              />
              {games.length > 1 && (
                <button
                  className="rounded-md border border-red-200 bg-red-50 px-2 py-2 text-xs font-semibold text-red-700"
                  onClick={() => removeGame(game.gameNumber)}
                  type="button"
                >
                  Remove
                </button>
              )}
            </div>
          ))}
        </div>

        {error && <p className="mb-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>}

        <div className="flex justify-end gap-3">
          <button
            className="rounded-lg border border-slate-300 bg-white px-4 py-2 text-sm font-semibold text-slate-700"
            onClick={onClose}
            type="button"
          >
            Cancel
          </button>
          <button
            className="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white disabled:opacity-50"
            disabled={submitting}
            onClick={() => void submit()}
            type="button"
          >
            {submitting ? "Submitting..." : "Submit result"}
          </button>
        </div>
      </div>
    </div>
  );
}

export default ResultDialog;
