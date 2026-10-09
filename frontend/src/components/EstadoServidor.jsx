export default function EstadoServidor({ servidor, onReintentar }) {
  if (servidor.estado === "listo") return null;
  if (servidor.estado === "iniciando") {
    return <p className="muted" role="status">Iniciando servidor, puede demorar unos segundos…</p>;
  }
  return <div className="auth-form">
    <p className="form-error" role="alert">{servidor.error}</p>
    <button className="secondary-button" type="button" onClick={onReintentar}>Volver a intentar</button>
  </div>;
}
