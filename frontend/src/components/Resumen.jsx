import { formatearMoneda } from "../utils/formato";

export default function Resumen({ resumen, cargando }) {
  return (
    <section className="section" aria-labelledby="resumen-title">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Vista general</p>
          <h2 id="resumen-title">Resumen</h2>
        </div>
      </div>
      {cargando ? (
        <p className="muted">Cargando resumen...</p>
      ) : (
        <>
          <div className="summary-grid">
            <article className="summary-card">
              <span>Total gastado</span>
              <strong>{formatearMoneda(resumen?.total || 0)}</strong>
            </article>
            <article className="summary-card">
              <span>Cantidad de gastos</span>
              <strong>{resumen?.cantidadGastos || 0}</strong>
            </article>
          </div>
          {resumen?.porCategoria?.length > 0 && (
            <div className="category-summary" aria-label="Gastos por categoría">
              {resumen.porCategoria.map((item) => (
                <span key={item.categoriaId}>{item.categoria}: <b>{formatearMoneda(item.total)}</b></span>
              ))}
            </div>
          )}
        </>
      )}
    </section>
  );
}
