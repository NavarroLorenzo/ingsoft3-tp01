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

export function validarFechaGasto(fecha) {
  if (typeof fecha !== "string" || fecha.trim() === "") {
    return "La fecha es obligatoria.";
  }

  const [anio, mes, dia] = fecha.split("-").map(Number);
  const fechaParseada = new Date(`${fecha}T00:00:00Z`);
  const formatoISO = /^\d{4}-\d{2}-\d{2}$/.test(fecha);
  const fechaExiste = !Number.isNaN(fechaParseada.getTime())
    && fechaParseada.getUTCFullYear() === anio
    && fechaParseada.getUTCMonth() + 1 === mes
    && fechaParseada.getUTCDate() === dia;

  if (!formatoISO || !fechaExiste) {
    return "La fecha debe tener el formato YYYY-MM-DD y ser una fecha válida.";
  }

  return null;
}

export function validarDescripcionGasto(descripcion) {
  if (typeof descripcion !== "string") return "La descripción es obligatoria.";

  const texto = descripcion.trim();
  if (texto === "") return "La descripción es obligatoria.";
  if (texto.length < 3) return "La descripción debe tener al menos 3 caracteres.";
  if (texto.length > 200) return "La descripción no puede superar los 200 caracteres.";
  return null;
}

export function validarGasto(gasto) {
  const descripcion = gasto.descripcion.trim();
  if (descripcion.length < 3 || descripcion.length > 200) return "La descripción debe tener entre 3 y 200 caracteres.";
  const monto = Number(gasto.monto);
  if (!Number.isFinite(monto) || monto <= 0) {
    return "El monto debe ser mayor que cero.";
  }
  if (Math.abs(monto * 100 - Math.round(monto * 100)) > 0.000001) {
    return "El monto puede tener como máximo dos decimales.";
  }
  const errorFecha = validarFechaGasto(gasto.fecha);
  if (errorFecha) return errorFecha;
  if (!gasto.categoriaId) return "La categoría es obligatoria.";
  return null;
}

export function validarCategoria(nombre) {
  const length = nombre.trim().length;
  return length < 2 || length > 50 ? "El nombre debe tener entre 2 y 50 caracteres." : null;
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
