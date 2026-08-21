export default function Header({ user, onLogout }) {
  return (
    <header className="site-header">
      <div>
        <p className="eyebrow">Control simple de finanzas</p>
        <h1>Gestor de Gastos Personales</h1>
      </div>
      <div className="user-actions">
        <span>Hola, {user?.nombre}</span>
        <button className="secondary-button" type="button" onClick={onLogout}>Cerrar sesión</button>
      </div>
    </header>
  );
}
