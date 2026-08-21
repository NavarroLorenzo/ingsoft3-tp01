export function formatearMoneda(monto) {
  return new Intl.NumberFormat("es-AR", {
    style: "currency",
    currency: "ARS"
  }).format(Number(monto));
}

export function formatearFecha(fecha) {
  if (!fecha) return "";
  return new Intl.DateTimeFormat("es-AR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric"
  }).format(new Date(`${fecha}T00:00:00`));
}

export function fechaActual() {
  return new Date().toLocaleDateString("en-CA");
}

export function validarGasto(gasto) {
  const descripcion = gasto.descripcion.trim();
  if (descripcion.length < 3 || descripcion.length > 200) {
    return "La descripción debe tener entre 3 y 200 caracteres.";
  }
  const monto = Number(gasto.monto);
  if (!Number.isFinite(monto) || monto <= 0) {
    return "El monto debe ser mayor que cero.";
  }
  if (Math.abs(monto * 100 - Math.round(monto * 100)) > 0.000001) {
    return "El monto puede tener como máximo dos decimales.";
  }
  if (!gasto.fecha) return "La fecha es obligatoria.";
  if (!gasto.categoriaId) return "La categoría es obligatoria.";
  return null;
}

export function validarEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim());
}

export function validarRegistro(datos) {
  if (datos.nombre.trim().length < 2 || datos.nombre.trim().length > 100) return "El nombre debe tener entre 2 y 100 caracteres.";
  if (!validarEmail(datos.email)) return "Ingresá un email válido.";
  if (datos.password.length < 8 || datos.password.length > 72) return "La contraseña debe tener entre 8 y 72 caracteres.";
  if (datos.password !== datos.confirmacion) return "Las contraseñas no coinciden.";
  return null;
}
