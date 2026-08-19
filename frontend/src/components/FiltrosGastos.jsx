import { useState } from "react";

const emptyFilters = { categoriaId: "", desde: "", hasta: "", texto: "" };

export default function FiltrosGastos({ categorias, onFiltrar, cargando }) {
  const [filtros, setFiltros] = useState(emptyFilters);

  function submit(event) {
    event.preventDefault();
    onFiltrar(filtros);
  }

  function limpiar() {
    setFiltros(emptyFilters);
    onFiltrar(emptyFilters);
  }

  return (
    <section className="section" aria-labelledby="filtros-title">
      <div className="section-heading"><h2 id="filtros-title">Filtros</h2></div>
      <form className="filters-form" onSubmit={submit}>
        <label>Categoría
          <select value={filtros.categoriaId} onChange={(event) => setFiltros({ ...filtros, categoriaId: event.target.value })}>
            <option value="">Todas</option>
            {categorias.map((categoria) => <option key={categoria.id} value={categoria.id}>{categoria.nombre}</option>)}
          </select>
        </label>
        <label>Desde
          <input type="date" value={filtros.desde} onChange={(event) => setFiltros({ ...filtros, desde: event.target.value })} />
        </label>
        <label>Hasta
          <input type="date" value={filtros.hasta} onChange={(event) => setFiltros({ ...filtros, hasta: event.target.value })} />
        </label>
        <label>Buscar
          <input type="search" value={filtros.texto} placeholder="Descripción" onChange={(event) => setFiltros({ ...filtros, texto: event.target.value })} />
        </label>
        <div className="form-actions filter-actions">
          <button className="primary-button" type="submit" disabled={cargando}>Filtrar</button>
          <button className="secondary-button" type="button" onClick={limpiar} disabled={cargando}>Limpiar</button>
        </div>
      </form>
    </section>
  );
}
