# ShuttleHub Badminton Tournament

A full-stack badminton tournament registration portal for managing men's doubles events. Players can view event availability, submit a team registration, and browse confirmed teams. Administrators can securely review registrations, confirm or reject teams, and withdraw teams while preserving registration history.

## Features

### Public portal

- Tournament landing page with live event status and capacity
- Men's doubles team registration
- Indian mobile number validation
- Required Player 1 phone number
- Optional Player 2 phone number
- Optional team name with player-name fallback
- Confirmed teams listing
- ShuttleHub branding and linked homepage navigation

### Admin dashboard

- JWT-based administrator login
- Protected registration dashboard
- Registration summary by status
- Confirm or reject pending registrations
- Withdraw teams without permanently deleting registration history
- Automatic redirect to login when authentication expires or is invalid

### Data integrity

- Registration status and capacity rules enforced by the backend
- Withdrawn and rejected teams do not consume capacity
- Duplicate participant protection per event
- Unique supplied phone numbers at the database level
- Database transactions for registration creation
- PostgreSQL foreign keys and cascade rules for tournament ownership

## Technology Stack

### Frontend

- React 19
- TypeScript
- Vite
- React Router DOM
- Tailwind CSS
- Oxlint

### Backend

- Go
- Chi router
- PostgreSQL
- pgx connection pool
- JWT authentication
- bcrypt password hashing
- golang-migrate-compatible SQL migrations

## Project Structure

```text
badminton-tournament/
├── backend/
│   ├── cmd/
│   │   ├── create-admin/       # CLI for creating an administrator
│   │   └── server/             # HTTP API server
│   ├── internal/
│   │   ├── database/           # PostgreSQL connection setup
│   │   ├── handler/             # HTTP handlers
│   │   ├── middleware/          # Authentication middleware
│   │   ├── model/               # Request and response models
│   │   ├── repository/          # Database queries
│   │   ├── router/              # API route registration
│   │   └── service/             # Business rules
│   ├── migrations/              # SQL up/down migrations
│   ├── go.mod
│   └── .env                    # Local only, never commit
├── frontend/
│   ├── public/                  # Static assets and ShuttleHub logo
│   ├── src/
│   │   ├── api/                 # Backend API clients
│   │   ├── components/          # Shared React components
│   │   ├── pages/               # Public and admin pages
│   │   └── types/               # Frontend types
│   ├── package.json
│   └── .env                    # Local only, never commit
└── README.md
```

## Application Routes

| Route | Description | Access |
| --- | --- | --- |
| `/` | Tournament homepage | Public |
| `/register` | Register a men's doubles team | Public |
| `/teams` | Browse confirmed teams | Public |
| `/admin/login` | Administrator login | Public |
| `/admin/registrations` | Manage registrations | Admin |

## API Endpoints

| Method | Endpoint | Description | Access |
| --- | --- | --- | --- |
| `GET` | `/health` | Health check | Public |
| `GET` | `/api/v1/events/{eventID}/` | Get event details | Public |
| `POST` | `/api/v1/events/{eventID}/registrations` | Register a team | Public |
| `GET` | `/api/v1/events/{eventID}/teams` | List confirmed teams | Public |
| `POST` | `/api/v1/auth/login` | Sign in as administrator | Public |
| `GET` | `/api/v1/admin/me` | Get authenticated administrator | Admin |
| `GET` | `/api/v1/admin/events/{eventID}/registrations` | List registrations | Admin |
| `PUT` | `/api/v1/admin/registrations/{registrationID}/status` | Update status | Admin |
| `PUT` | `/api/v1/admin/registrations/{registrationID}/withdraw` | Withdraw a registration | Admin |

Admin endpoints require an authorization header:

```http
Authorization: Bearer <jwt-token>
```

## Prerequisites

Install the following before running the project:

- Go 1.27 or compatible
- Node.js and npm
- PostgreSQL
- `migrate` CLI compatible with the migration files

The backend expects a PostgreSQL database named `badminton_tournament` unless you choose another database in `DATABASE_URL`.

## Database Setup

Create the database in PostgreSQL, then configure the backend environment:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/badminton_tournament?sslmode=disable
PORT=8080
JWT_SECRET=replace-with-a-long-random-secret
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

Apply migrations from the backend directory:

```bash
cd backend
migrate -path migrations -database "$DATABASE_URL" up
```

Check the migration version:

```bash
migrate -path migrations -database "$DATABASE_URL" version
```

Do not commit `.env` files or real secrets. Use a strong random `JWT_SECRET` outside local development.

## Create an Administrator

Set an administrator password in the backend environment:

```env
ADMIN_PASSWORD=choose-a-local-admin-password
```

Run the admin creation command from the backend directory:

```bash
go run ./cmd/create-admin
```

The current development account email created by this command is:

```text
admin@badminton.local
```

Change this behavior before production if administrator email management is required.

## Run Locally

### Start the backend

```bash
cd backend
go run ./cmd/server
```

The API runs at:

```text
http://localhost:8080
```

### Start the frontend

Create `frontend/.env`:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Then run:

```bash
cd frontend
npm install
npm run dev
```

The frontend runs at:

```text
http://localhost:5173
```

If port `8080` is already in use, either stop the existing backend process or run the backend on another port:

```powershell
$env:PORT = "8081"
go run ./cmd/server
```

Update `VITE_API_BASE_URL` to match the new port.

## Validation Rules

- Player 1 name is required.
- Player 1 phone is required.
- Player 2 name is required.
- Player 2 phone is optional.
- When supplied, Indian mobile numbers must contain exactly 10 digits and start with `6`, `7`, `8`, or `9`.
- Player 1 and Player 2 cannot use the same phone number.
- A participant cannot register twice in the same event when their phone is available.
- Registration is rejected when the event is closed, past its deadline, or full.

The mobile-number format check confirms a basic Indian numbering pattern. It does not prove that a number is active or belongs to the participant. OTP verification would be needed for that.

## Development Commands

### Frontend

```bash
cd frontend
npm run dev       # Start Vite development server
npm run lint      # Run Oxlint
npm run build     # Type-check and create production build
npm run preview   # Preview the production build
```

### Backend

```bash
cd backend
go test ./...     # Run all Go tests
go test -race ./... # Run tests with race detection
go vet ./...      # Run Go static analysis
gofmt -w .        # Format Go source files
```

## Testing

The backend tests are under `backend/tests`. They cover registration validation and registration error-to-HTTP-status mapping.

The test suite verifies, among other cases:

- Optional Player 2 phone numbers
- Invalid Indian phone formats
- Duplicate phone numbers
- Missing required fields
- Event-not-found responses
- Registration conflicts
- Unexpected server errors

Before opening a pull request, run:

```bash
cd backend
go test -race ./...
go vet ./...

cd ../frontend
npm run lint
npm run build
```

## Status Handling

The main registration statuses are:

- `PENDING`: submitted and awaiting admin review
- `CONFIRMED`: visible on the public teams page
- `REJECTED`: not visible publicly and does not consume capacity
- `WITHDRAWN`: removed from the public teams page, does not consume capacity, and remains in the admin history

Withdrawal is intentionally preferred over hard deletion for published tournaments because it preserves an audit trail.

## Deployment Notes

Before deploying:

1. Use production PostgreSQL credentials and a strong JWT secret.
2. Apply all migrations in order.
3. Set `CORS_ALLOWED_ORIGINS` to the exact deployed frontend origin.
4. Set the frontend `VITE_API_BASE_URL` to the deployed API origin.
5. Serve the frontend `dist` directory through a static host or web server.
6. Configure SPA fallback to `index.html` so direct navigation to `/teams` and admin routes works.
7. Use HTTPS in production.
8. Confirm database backups and migration rollback procedures.
9. Do not expose administrator credentials in source control or deployment logs.

## License

This project is intended for internal tournament management and demonstration purposes.
