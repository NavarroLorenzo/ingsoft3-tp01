import { useEffect, useState } from "react";
import {
  actualizarCategoria, actualizarGasto, crearCategoria, crearGasto, eliminarCategoria,
  eliminarGasto, obtenerCategorias, obtenerGastos, obtenerResumen
} from "./api/api";
import Header from "./components/Header";
import Resumen from "./components/Resumen";
import FormularioGasto from "./components/FormularioGasto";
import FiltrosGastos from "./components/FiltrosGastos";
import ListaGastos from "./components/ListaGastos";
import Categorias from "./components/Categorias";

const initialFilters = { categoriaId: "", desde: "", hasta: "", texto: "" };

export default function Dashboard({ user, onLogout }) {
  const [categorias, setCategorias] = useState([]);
  const [gastos, setGastos] = useState([]);
  const [resumen, setResumen] = useState(null);
  const [filtros, setFiltros] = useState(initialFilters);
  const [gastoEdicion, setGastoEdicion] = useState(null);
  const [cargando, setCargando] = useState(true);
  const [guardando, setGuardando] = useState(false);
  const [error, setError] = useState("");
  const [exito, setExito] = useState("");

  useEffect(() => { cargarTodo(initialFilters); }, []);

  async function cargarTodo(currentFilters = filtros) {
    setCargando(true);
    setError("");
    try {
      const [nuevasCategorias, nuevosGastos, nuevoResumen] = await Promise.all([
        obtenerCategorias(), obtenerGastos(currentFilters), obtenerResumen({ desde: currentFilters.desde, hasta: currentFilters.hasta })
      ]);
      setCategorias(nuevasCategorias);
      setGastos(nuevosGastos);
      setResumen(nuevoResumen);
    } catch (requestError) {
      setError(requestError.message || "No se pudieron obtener los datos.");
    } finally {
      setCargando(false);
    }
  }

  function mostrarExito(message) {
    setExito(message);
    window.setTimeout(() => setExito(""), 3500);
  }

  async function guardarGasto(gasto) {
    setGuardando(true);
    setError("");
    try {
      if (gastoEdicion) {
        await actualizarGasto(gastoEdicion.id, gasto);
        setGastoEdicion(null);
        mostrarExito("Gasto actualizado correctamente.");
      } else {
        await crearGasto(gasto);
        mostrarExito("Gasto registrado correctamente.");
      }
      await cargarTodo();
      return true;
    } catch (requestError) {
      setError(requestError.message);
      return false;
    } finally {
      setGuardando(false);
    }
  }

  async function filtrarGastos(nuevosFiltros) {
    setFiltros(nuevosFiltros);
    await cargarTodo(nuevosFiltros);
  }

  async function borrarGasto(gasto) {
    if (!window.confirm(`¿Eliminar el gasto “${gasto.descripcion}”?`)) return;
    try {
      await eliminarGasto(gasto.id);
      if (gastoEdicion?.id === gasto.id) setGastoEdicion(null);
      mostrarExito("Gasto eliminado correctamente.");
      await cargarTodo();
    } catch (requestError) { setError(requestError.message); }
  }

  async function performCategory(action, successMessage) {
    setError("");
    try {
      await action();
      mostrarExito(successMessage);
      await cargarTodo();
      return true;
    } catch (requestError) {
      setError(requestError.message);
      return false;
    }
  }

  function borrarCategoria(categoria) {
    if (!window.confirm(`¿Eliminar la categoría “${categoria.nombre}”?`)) return;
    performCategory(() => eliminarCategoria(categoria.id), "Categoría eliminada correctamente.");
  }

  return (
    <main className="app-shell">
      <Header user={user} onLogout={onLogout} />
      {error && <div className="notification error" role="alert">{error}</div>}
      {exito && <div className="notification success" role="status">{exito}</div>}
      <Resumen resumen={resumen} cargando={cargando} />
      <div className="two-columns">
        <FormularioGasto categorias={categorias} gastoEdicion={gastoEdicion} onGuardar={guardarGasto} onCancelar={() => setGastoEdicion(null)} cargando={guardando} />
        <FiltrosGastos categorias={categorias} onFiltrar={filtrarGastos} cargando={cargando} />
      </div>
      <ListaGastos gastos={gastos} cargando={cargando} onEditar={setGastoEdicion} onEliminar={borrarGasto} />
      <Categorias categorias={categorias} cargando={guardando} onCrear={(categoria) => performCategory(() => crearCategoria(categoria), "Categoría creada correctamente.")} onActualizar={(id, categoria) => performCategory(() => actualizarCategoria(id, categoria), "Categoría actualizada correctamente.")} onEliminar={borrarCategoria} />
    </main>
  );
}
