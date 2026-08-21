let onUnauthorized = null;

export function setUnauthorizedHandler(handler) {
  onUnauthorized = handler;
}

export function getToken() {
  return localStorage.getItem("token");
}

export function guardarSesion({ token, user }) {
  localStorage.setItem("token", token);
  localStorage.setItem("user", JSON.stringify(user));
}

export function limpiarSesion() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}

export function getAuthHeaders() {
  const token = getToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request(path, options = {}) {
  const { auth = true, ...fetchOptions } = options;
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...(auth ? getAuthHeaders() : {}),
      ...fetchOptions.headers
    },
    ...fetchOptions
  });

  if (response.status === 204) return null;
  const body = await response.json().catch(() => null);
  if (!response.ok) {
    if (response.status === 401 && auth) {
      limpiarSesion();
      onUnauthorized?.();
    }
    throw new Error(body?.error || "Ocurrió un error al comunicarse con el servidor.");
  }
  return body;
}

function queryString(filtros = {}) {
  const parameters = new URLSearchParams();
  Object.entries(filtros).forEach(([key, value]) => {
    if (value !== "" && value !== undefined && value !== null) parameters.set(key, value);
  });
  const query = parameters.toString();
  return query ? `?${query}` : "";
}

export function register(datos) {
  return request("/api/auth/register", { auth: false, method: "POST", body: JSON.stringify(datos) });
}

export function login(datos) {
  return request("/api/auth/login", { auth: false, method: "POST", body: JSON.stringify(datos) });
}

export function obtenerUsuarioActual() {
  return request("/api/auth/me");
}

export function obtenerGastos(filtros) { return request(`/api/gastos${queryString(filtros)}`); }
export function crearGasto(gasto) { return request("/api/gastos", { method: "POST", body: JSON.stringify(gasto) }); }
export function actualizarGasto(id, gasto) { return request(`/api/gastos/${id}`, { method: "PUT", body: JSON.stringify(gasto) }); }
export function eliminarGasto(id) { return request(`/api/gastos/${id}`, { method: "DELETE" }); }
export function obtenerCategorias() { return request("/api/categorias"); }
export function crearCategoria(categoria) { return request("/api/categorias", { method: "POST", body: JSON.stringify(categoria) }); }
export function actualizarCategoria(id, categoria) { return request(`/api/categorias/${id}`, { method: "PUT", body: JSON.stringify(categoria) }); }
export function eliminarCategoria(id) { return request(`/api/categorias/${id}`, { method: "DELETE" }); }
export function obtenerResumen(filtros) { return request(`/api/resumen${queryString(filtros)}`); }
