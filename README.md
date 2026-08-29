# ShuttleHub Badminton Tournament

ShuttleHub is a full-stack badminton tournament platform for managing events, team registrations, draws, round progression, and live match results. The project combines a React and TypeScript frontend with a Go backend powered by PostgreSQL, supporting both public tournament participation and admin-side competition control.

## What the project does

### Public experience

- Browse tournament events and event details
- View live event status, registration capacity, and venue information
- Register a men’s doubles team with validation for Indian phone numbers
- See confirmed teams on the public teams page
- Open event-specific routes for registration, team listing, and match details

### Admin experience

- Sign in securely as an admin or super admin
- Review and manage registrations from an admin dashboard
- Confirm, reject, or withdraw registrations while preserving history
- Create and manage tournament rounds and bracket draws
- Enter match results and progress the bracket automatically
- Handle three-team semifinal workflows and auto-create downstream fixtures
- Manage platform-level events and admin assignments from the super-admin dashboard

### Recent platform capabilities

- Responsive shared header navigation for desktop and mobile
- Secure authentication with HTTP-only cookies and role-based access
- CSRF origin validation for state-changing admin requests
- Rate limiting for login and registration endpoints
- Automatic bracket progression for byes and playoff scenarios
- Transaction-safe draw generation and result submission

## Technology stack

### Frontend

- React 19
- TypeScript
- Vite
- React Router
- Tailwind CSS
- Oxlint

### Backend

- Go 1.27
- Chi router
- PostgreSQL
- pgx connection pool
- JWT authentication
- bcrypt password hashing
- SQL migrations via golang-migrate-compatible files

## Project structure

```text
badminton-tournament/
├── backend/
│   ├── cmd/
│   │   ├── create-admin/
│   │   └── server/
│   ├── internal/
│   │   ├── database/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── router/
│   │   └── service/
│   ├── migrations/
│   └── go.mod
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── pages/
│   │   └── types/
│   └── package.json
└── README.md
```

## Main routes

- `/` — public home page
- `/events/:eventId` — event details and status
- `/events/:eventId/register` — team registration
- `/events/:eventId/teams` — confirmed teams
- `/events/:eventId/live-scores` — match and round results
- `/admin/login` — admin login
- `/admin` — admin control center
- `/admin/registrations` — registration review
- `/admin/draw` — round and draw management
- `/admin/superadmin` — super-admin dashboard

## Prerequisites

Install the following before running the app locally:

- Go 1.27 or newer
- Node.js and npm
- PostgreSQL

## Backend environment

Create a backend environment file with values such as:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/badminton_tournament?sslmode=disable
PORT=8080
JWT_SECRET=replace-with-a-long-random-secret
COOKIE_SECURE=false
CORS_ALLOWED_ORIGINS=http://localhost:5173
LOGIN_RATE_LIMIT_MAX_REQUESTS=5
LOGIN_RATE_LIMIT_WINDOW_SECONDS=60
REGISTRATION_RATE_LIMIT_MAX_REQUESTS=10
REGISTRATION_RATE_LIMIT_WINDOW_SECONDS=60
```

Apply migrations from the backend directory:

```bash
cd backend
migrate -path migrations -database "$DATABASE_URL" up
```

Create an initial admin account:

```bash
go run ./cmd/create-admin
```

## Run locally

### Start the backend

```bash
cd backend
go run ./cmd/server
```

The API will run on:

```text
http://localhost:8080
```

### Start the frontend

Create a frontend environment file:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Then run:

```bash
cd frontend
npm install
npm run dev
```

The frontend will run on:

```text
http://localhost:5173
```

## Validation and business rules

- Player 1 name is required
- Player 1 phone is required
- Player 2 name is optional, but if provided it must follow the same validation rules
- Indian mobile numbers must be 10 digits and start with 6, 7, 8, or 9
- Duplicate phone numbers are blocked within an event
- Registrations are rejected when the event is closed, past its deadline, or full
- Registration state transitions are enforced server-side
- Draw generation and match results are handled transactionally

## Development commands

### Frontend

```bash
cd frontend
npm run dev
npm run lint
npm run build
npm run preview
```

### Backend

```bash
cd backend
go test ./...
go vet ./...
gofmt -w .
```

## Testing

The backend test suite covers core behavior such as registration validation, state transitions, draw generation, and result handling.

Before opening a pull request, it is recommended to run:

```bash
cd backend
go test ./...
go vet ./...

cd ../frontend
npm run lint
npm run build
```

## Notes

This project is intended for internal tournament management and demonstration purposes. It currently supports event-based tournament workflows with a public registration flow and a secure admin console for managing competition operations.
