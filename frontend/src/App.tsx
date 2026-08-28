import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";

import Home from "./pages/Home";
import Registration from "./pages/Registration";
import Teams from "./pages/Teams";
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
          path="/teams"
          element={<Teams />}
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
          path="/admin/registrations"
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

      </Routes>
    </BrowserRouter>
  );
}

export default App;