# Badminton Tournament Frontend

A modern React + TypeScript frontend for the badminton tournament registration system. It provides a public-facing event experience for players and a protected admin dashboard for managing team registrations.

## Overview

This frontend app lets users:

- view the tournament home page with current event status
- register a men's doubles team with player details
- see the list of registered teams and their status
- sign in as an admin to review registrations
- confirm or reject pending registrations from a secure admin panel

## Features

### Public tournament experience

- Landing page for the tournament
- Live event details such as registration status and team capacity
- Team registration form for two players and an optional team name
- Player 2 phone number is optional; player 1 phone remains required
- Phone numbers must contain exactly 10 digits and start with 6, 7, 8, or 9
- Success state after a registration is submitted
- Team list page to browse all registered entries

### Admin features

- Admin sign-in page with token-based authentication
- Protected route guarding for admin-only pages
- Registration dashboard showing player information and team names
- Status controls to confirm or reject pending entries
- Withdrawal controls preserve registration history while removing teams from the public list
- Automatic sign-out on unauthorized access

### Tech stack

- React 19
- TypeScript
- Vite
- React Router DOM
- Tailwind CSS

## Project structure

```bash
frontend/
├── src/
│   ├── api/              # API calls for auth, events, registrations, and teams
│   ├── components/       # Shared UI components such as protected routes
│   ├── pages/            # Public and admin pages
│   ├── types/            # Shared TypeScript types
│   ├── App.tsx           # Route configuration
│   ├── main.tsx          # Application entry point
│   └── index.css         # Global styles
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## Available routes

- `/` — Tournament home page
- `/register` — Team registration form
- `/teams` — Registered teams overview
- `/admin/login` — Administrator login
- `/admin/registrations` — Protected registration management dashboard

## Setup

1. Install dependencies:

```bash
npm install
```

2. Create a `.env` file in the `frontend` folder and configure the backend API URL:

```bash
VITE_API_BASE_URL=http://localhost:8080
```

> Adjust the value to match your backend host and port.

The backend also needs `DATABASE_URL`, `JWT_SECRET`, and (for deployed
frontends) `CORS_ALLOWED_ORIGINS` configured in its environment.

3. Start the development server:

```bash
npm run dev
```

The app will run in Vite development mode, typically on:

```bash
http://localhost:5173
```

## Production build

```bash
npm run build
```

This runs TypeScript checks and produces a production build in the `dist` folder.

## Scripts

```bash
npm run dev      # run the development server
npm run build    # compile and build the app
npm run preview  # preview the production build locally
npm run lint     # run the project linter
```

## Notes

- The frontend currently targets a fixed tournament event ID for the registration flow.
- A supplied phone number is unique in the database, and team player IDs are unique per event. A phone-less player cannot be reliably deduplicated by name alone; for stronger identity protection, require either a phone number or another stable identifier such as email.
- Admin authentication is stored in localStorage using a JWT-like token.
- Protected admin pages redirect users back to the login page when the session is invalid or expired.

## License

This project is for internal tournament management and demonstration purposes.
