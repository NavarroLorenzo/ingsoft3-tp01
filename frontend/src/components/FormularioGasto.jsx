import { useEffect, useState } from "react";
import { fechaActual, validarGasto } from "../utils/formato";

const initialGasto = () => ({ descripcion: "", monto: "", fecha: fechaActual(), categoriaId: "" });

export default function FormularioGasto({ categorias, gastoEdicion, onGuardar, onCancelar, cargando }) {
  const [gasto, setGasto] = useState(initialGasto);
  const [error, setError] = useState("");

  useEffect(() => {
    setError("");
    if (gastoEdicion) {
      setGasto({
        descripcion: gastoEdicion.descripcion,
        monto: gastoEdicion.monto,
        fecha: gastoEdicion.fecha,
        categoriaId: String(gastoEdicion.categoriaId)
      });
    } else {
      setGasto(initialGasto());
    }
  }, [gastoEdicion]);

  function change(event) {
    setGasto((current) => ({ ...current, [event.target.name]: event.target.value }));
  }

  async function submit(event) {
    event.preventDefault();
    const validationError = validarGasto(gasto);
    if (validationError) {
      setError(validationError);
      return;
    }
    setError("");
    const saved = await onGuardar({ ...gasto, monto: Number(gasto.monto), categoriaId: Number(gasto.categoriaId) });
    if (saved) setGasto(initialGasto());
  }

  return (
    <section className="section" aria-labelledby="gasto-form-title">
      <div className="section-heading"><h2 id="gasto-form-title">{gastoEdicion ? "Editar gasto" : "Nuevo gasto"}</h2></div>
      <form className="expense-form" onSubmit={submit}>
        <label>Descripción
          <input name="descripcion" type="text" value={gasto.descripcion} onChange={change} minLength="3" maxLength="200" required />
        </label>
        <label>Monto
          <input name="monto" type="number" value={gasto.monto} onChange={change} min="0.01" step="0.01" required />
        </label>
        <label>Fecha
          <input name="fecha" type="date" value={gasto.fecha} onChange={change} required />
        </label>
        <label>Categoría
          <select name="categoriaId" value={gasto.categoriaId} onChange={change} required>
            <option value="">Seleccionar categoría</option>
            {categorias.map((categoria) => <option key={categoria.id} value={categoria.id}>{categoria.nombre}</option>)}
          </select>
        </label>
        {error && <p className="form-error" role="alert">{error}</p>}
        <div className="form-actions">
          <button className="primary-button" type="submit" disabled={cargando || categorias.length === 0}>{cargando ? "Guardando..." : gastoEdicion ? "Guardar cambios" : "Registrar gasto"}</button>
          {gastoEdicion && <button className="secondary-button" type="button" onClick={onCancelar}>Cancelar</button>}
        </div>
      </form>
    </section>
  );
}
