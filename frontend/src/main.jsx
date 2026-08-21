import { StrictMode, useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import App from "./App";
import ProtectedRoute from "./components/ProtectedRoute";
import { LoginPage, RegisterPage } from "./components/AuthForm";
import { limpiarSesion, guardarSesion, getToken, login, obtenerUsuarioActual, register, setUnauthorizedHandler } from "./api/api";
import "./styles.css";

function RouterApp() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    setUnauthorizedHandler(() => { setUser(null); navigate("/login", { replace: true }); });
    if (!getToken()) { setLoading(false); return; }
    obtenerUsuarioActual().then(setUser).catch(() => { limpiarSesion(); setUser(null); }).finally(() => setLoading(false));
  }, [navigate]);

  async function startSession(action, datos) {
    const session = await action(datos);
    guardarSesion(session);
    setUser(session.user);
    navigate("/dashboard", { replace: true });
  }

  function logout() { limpiarSesion(); setUser(null); navigate("/login", { replace: true }); }
  if (loading) return <main className="auth-page"><p className="muted">Validando sesión...</p></main>;

  return <Routes>
    <Route path="/login" element={user ? <Navigate to="/dashboard" replace /> : <LoginPage onLogin={(datos) => startSession(login, datos)} />} />
    <Route path="/register" element={user ? <Navigate to="/dashboard" replace /> : <RegisterPage onRegister={(datos) => startSession(register, datos)} />} />
    <Route element={<ProtectedRoute authenticated={Boolean(user)} />}>
      <Route path="/dashboard" element={<App user={user} onLogout={logout} />} />
    </Route>
    <Route path="*" element={<Navigate to={user ? "/dashboard" : "/login"} replace />} />
  </Routes>;
}

createRoot(document.getElementById("root")).render(<StrictMode><BrowserRouter><RouterApp /></BrowserRouter></StrictMode>);
