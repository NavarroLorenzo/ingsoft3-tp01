import { Link } from "react-router-dom";
import { useState } from "react";
import { validarEmail, validarRegistro } from "../utils/formato";

export function LoginPage({ onLogin }) {
  const [datos, setDatos] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  async function submit(event) {
    event.preventDefault();
    if (!validarEmail(datos.email) || !datos.password) {
      setError("Ingresá un email válido y tu contraseña.");
      return;
    }
    setCargando(true); setError("");
    try { await onLogin(datos); } catch (requestError) { setError(requestError.message); } finally { setCargando(false); }
  }
  return <AuthCard title="Iniciar sesión" onSubmit={submit} error={error} loading={cargando} button="Iniciar sesión">
    <label>Email<input type="email" value={datos.email} onChange={(e) => setDatos({ ...datos, email: e.target.value })} required autoComplete="email" /></label>
    <label>Contraseña<input type="password" value={datos.password} onChange={(e) => setDatos({ ...datos, password: e.target.value })} required autoComplete="current-password" /></label>
    <p className="auth-link">¿No tenés cuenta? <Link to="/register">Registrate</Link></p>
  </AuthCard>;
}

export function RegisterPage({ onRegister }) {
  const [datos, setDatos] = useState({ nombre: "", email: "", password: "", confirmacion: "" });
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);
  async function submit(event) {
    event.preventDefault();
    const validationError = validarRegistro(datos);
    if (validationError) { setError(validationError); return; }
    setCargando(true); setError("");
    try { await onRegister({ nombre: datos.nombre, email: datos.email, password: datos.password }); } catch (requestError) { setError(requestError.message); } finally { setCargando(false); }
  }
  return <AuthCard title="Crear cuenta" onSubmit={submit} error={error} loading={cargando} button="Crear cuenta">
    <label>Nombre<input type="text" value={datos.nombre} onChange={(e) => setDatos({ ...datos, nombre: e.target.value })} required minLength="2" maxLength="100" autoComplete="name" /></label>
    <label>Email<input type="email" value={datos.email} onChange={(e) => setDatos({ ...datos, email: e.target.value })} required autoComplete="email" /></label>
    <label>Contraseña<input type="password" value={datos.password} onChange={(e) => setDatos({ ...datos, password: e.target.value })} required minLength="8" maxLength="72" autoComplete="new-password" /></label>
    <label>Confirmar contraseña<input type="password" value={datos.confirmacion} onChange={(e) => setDatos({ ...datos, confirmacion: e.target.value })} required minLength="8" maxLength="72" autoComplete="new-password" /></label>
    <p className="auth-link">¿Ya tenés cuenta? <Link to="/login">Iniciar sesión</Link></p>
  </AuthCard>;
}

function AuthCard({ title, onSubmit, error, loading, button, children }) {
  return <main className="auth-page"><section className="auth-card"><p className="eyebrow">Gestor de Gastos</p><h1>{title}</h1><form className="auth-form" onSubmit={onSubmit}>{children}{error && <p className="form-error" role="alert">{error}</p>}<button className="primary-button" disabled={loading} type="submit">{loading ? "Procesando..." : button}</button></form></section></main>;
}
