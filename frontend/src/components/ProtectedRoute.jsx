import { Navigate, Outlet } from "react-router-dom";

export default function ProtectedRoute({ authenticated }) {
  return authenticated ? <Outlet /> : <Navigate to="/login" replace />;
}
