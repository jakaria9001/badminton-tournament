import type { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";

import { getAdminToken } from "../api/authApi";

interface Props {
  children: ReactNode;
}

export default function ProtectedRoute({ children }: Props) {
  const location = useLocation();

  if (!getAdminToken()) {
    return (
      <Navigate
        replace
        state={{ from: location.pathname }}
        to="/admin/login"
      />
    );
  }

  return children;
}