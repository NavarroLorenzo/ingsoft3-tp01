import { StrictMode, useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import App from "./App";
import ProtectedRoute from "./components/ProtectedRoute";
import { LoginPage, RegisterPage } from "./components/AuthForm";
import { limpiarSesion, guardarSesion, getToken, login, obtenerUsuarioActual, register, setUnauthorizedHandler } from "./api/api";
import { iniciarSesion } from "./services/iniciarSesion";
import { esperarServidor } from "./services/esperarServidor";
import EstadoServidor from "./components/EstadoServidor";
import Brand from "./components/Brand";
import "./styles.css";

function RouterApp() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(() => Boolean(getToken()));
  const [servidor, setServidor] = useState({ estado: "iniciando", error: "" });
  const [intentoServidor, setIntentoServidor] = useState(0);
  const navigate = useNavigate();

  useEffect(() => {
    const controller = new AbortController();
    setServidor({ estado: "iniciando", error: "" });
    esperarServidor({ signal: controller.signal })
      .then(() => { if (!controller.signal.aborted) setServidor({ estado: "listo", error: "" }); })
      .catch((error) => { if (!controller.signal.aborted) setServidor({ estado: "error", error: error.message }); });
    return () => controller.abort();
  }, [intentoServidor]);

  useEffect(() => {
    setUnauthorizedHandler(() => { setUser(null); navigate("/login", { replace: true }); });
    if (servidor.estado !== "listo") return;
    if (!getToken()) { setLoading(false); return; }
    let activo = true;
    obtenerUsuarioActual()
      .then((usuario) => { if (activo) setUser(usuario); })
      .catch(() => { if (activo) { limpiarSesion(); setUser(null); } })
      .finally(() => { if (activo) setLoading(false); });
    return () => { activo = false; };
  }, [navigate, servidor.estado]);

  function reintentarServidor() {
    setServidor({ estado: "iniciando", error: "" });
    setIntentoServidor((intento) => intento + 1);
  }

  async function startSession(action, datos) {
    // También se comprueba al enviar: el formulario puede llevar tiempo abierto.
    setServidor({ estado: "iniciando", error: "" });
    try {
      await esperarServidor();
    } catch (error) {
      setServidor({ estado: "error", error: error.message });
      return;
    }
    setServidor({ estado: "listo", error: "" });
    const session = await iniciarSesion(datos, { autenticar: action, guardarSesion });
    setUser(session.user);
    navigate("/dashboard", { replace: true });
  }

  function logout() { limpiarSesion(); setUser(null); navigate("/login", { replace: true }); }
  if (loading) return <main className="auth-page"><section className="auth-card">
    <Brand className="auth-brand" />
    <h1>Ingresar a Cent</h1>
    <EstadoServidor servidor={servidor} onReintentar={reintentarServidor} />
    {servidor.estado === "listo" && <p className="muted" role="status">Validando sesión...</p>}
  </section></main>;

  return <Routes>
    <Route path="/login" element={user ? <Navigate to="/dashboard" replace /> : <LoginPage servidor={servidor} onReintentar={reintentarServidor} onLogin={(datos) => startSession(login, datos)} />} />
    <Route path="/register" element={user ? <Navigate to="/dashboard" replace /> : <RegisterPage servidor={servidor} onReintentar={reintentarServidor} onRegister={(datos) => startSession(register, datos)} />} />
    <Route element={<ProtectedRoute authenticated={Boolean(user)} />}>
      <Route path="/dashboard" element={<App user={user} onLogout={logout} />} />
    </Route>
    <Route path="*" element={<Navigate to={user ? "/dashboard" : "/login"} replace />} />
  </Routes>;
}

createRoot(document.getElementById("root")).render(<StrictMode><BrowserRouter><RouterApp /></BrowserRouter></StrictMode>);
