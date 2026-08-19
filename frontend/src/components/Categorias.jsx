import { useState } from "react";

export default function Categorias({ categorias, onCrear, onActualizar, onEliminar, cargando }) {
  const [nombre, setNombre] = useState("");
  const [edicion, setEdicion] = useState(null);
  const [error, setError] = useState("");

  async function submit(event) {
    event.preventDefault();
    const value = nombre.trim();
    if (value.length < 2 || value.length > 50) {
      setError("El nombre debe tener entre 2 y 50 caracteres.");
      return;
    }
    setError("");
    const result = edicion ? await onActualizar(edicion.id, { nombre: value }) : await onCrear({ nombre: value });
    if (result) {
      setNombre("");
      setEdicion(null);
    }
  }

  function startEdit(categoria) {
    setEdicion(categoria);
    setNombre(categoria.nombre);
    setError("");
  }

  return (
    <section className="section" aria-labelledby="categorias-title">
      <div className="section-heading"><h2 id="categorias-title">Administrar categorías</h2></div>
      <form className="category-form" onSubmit={submit}>
        <label className="sr-only" htmlFor="nombre-categoria">Nombre de categoría</label>
        <input id="nombre-categoria" type="text" value={nombre} maxLength="50" placeholder="Nueva categoría" onChange={(event) => setNombre(event.target.value)} required />
        <button className="primary-button" type="submit" disabled={cargando}>{edicion ? "Guardar" : "Agregar"}</button>
        {edicion && <button className="secondary-button" type="button" onClick={() => { setEdicion(null); setNombre(""); }}>Cancelar</button>}
      </form>
      {error && <p className="form-error" role="alert">{error}</p>}
      {categorias.length === 0 ? <p className="muted">No hay categorías registradas.</p> : (
        <ul className="category-list">
          {categorias.map((categoria) => (
            <li key={categoria.id}>
              <span>{categoria.nombre}</span>
              <span className="row-actions">
                <button className="text-button" type="button" onClick={() => startEdit(categoria)}>Editar</button>
                <button className="text-button danger" type="button" onClick={() => onEliminar(categoria)}>Eliminar</button>
              </span>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
