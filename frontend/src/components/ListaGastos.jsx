import FilaGasto from "./FilaGasto";

export default function ListaGastos({ gastos, cargando, onEditar, onEliminar }) {
  return (
    <section className="section" aria-labelledby="gastos-title">
      <div className="section-heading"><h2 id="gastos-title">Gastos</h2></div>
      {cargando ? <p className="muted">Cargando gastos...</p> : gastos.length === 0 ? <p className="muted">No hay gastos registrados.</p> : (
        <div className="table-scroll">
          <table>
            <thead><tr><th>Descripción</th><th>Categoría</th><th>Fecha</th><th>Monto</th><th>Acciones</th></tr></thead>
            <tbody>{gastos.map((gasto) => <FilaGasto key={gasto.id} gasto={gasto} onEditar={onEditar} onEliminar={onEliminar} />)}</tbody>
          </table>
        </div>
      )}
    </section>
  );
}
