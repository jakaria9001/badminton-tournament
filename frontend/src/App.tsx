import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";

import Home from "./pages/Home";
import EventDetails from "./pages/EventDetails";
import SuperAdminDashboard from "./pages/SuperAdminDashboard";
import Registration from "./pages/Registration";
import Teams from "./pages/Teams";
import LiveScores from "./pages/LiveScores";
import Registrations from "./pages/Admin/Registrations";
import Login from "./pages/Admin/Login";
import ControlCenter from "./pages/Admin/ControlCenter";
import Draw from "./pages/Admin/Draw";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <BrowserRouter>
      <Routes>

        <Route
          path="/"
          element={<Home />}
        />

        <Route
          path="/register"
          element={<Registration />}
        />

        <Route
          path="/events/:eventId/register"
          element={<Registration />}
        />

        <Route
          path="/teams"
          element={<Teams />}
        />

        <Route
          path="/events/:eventId/teams"
          element={<Teams />}
        />

        <Route
          path="/events/:eventId"
          element={<EventDetails />}
        />

        <Route
          path="/live-scores"
          element={<LiveScores />}
        />

        <Route
          path="/events/:eventId/live-scores"
          element={<LiveScores />}
        />

        <Route
          path="/admin/login"
          element={<Login />}
        />

        <Route
          path="/admin"
          element={
            <ProtectedRoute>
              <ControlCenter />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin/superadmin"
          element={
            <ProtectedRoute>
              <SuperAdminDashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin/registrations"
          element={
            <ProtectedRoute>
              <Registrations />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin/events/:eventId/registrations"
          element={
            <ProtectedRoute>
              <Registrations />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin/draw"
          element={
            <ProtectedRoute>
              <Draw />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin/events/:eventId/draw"
          element={
            <ProtectedRoute>
              <Draw />
            </ProtectedRoute>
          }
        />

      </Routes>
    </BrowserRouter>
  );
}

export default App;