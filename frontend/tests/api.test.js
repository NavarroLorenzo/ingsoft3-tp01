import { beforeEach, describe, expect, it, vi } from "vitest";
import {
  actualizarCategoria,
  actualizarGasto,
  crearCategoria,
  crearGasto,
  eliminarCategoria,
  eliminarGasto,
  guardarSesion,
  login,
  obtenerCategorias,
  obtenerGastos,
  obtenerResumen,
  obtenerUsuarioActual,
  register,
  setUnauthorizedHandler
} from "../src/api/api";

function respuesta(body, status = 200) {
  return {
    status,
    ok: status >= 200 && status < 300,
    json: vi.fn().mockResolvedValue(body)
  };
}

describe("cliente de API", () => {
  beforeEach(() => {
    const store = {};
    global.localStorage = {
      getItem: vi.fn((key) => store[key] || null),
      setItem: vi.fn((key, value) => { store[key] = value; }),
      removeItem: vi.fn((key) => { delete store[key]; })
    };
    global.fetch = vi.fn();
    setUnauthorizedHandler(null);
  });

  it.each([
    ["registra un usuario", () => register({ nombre: "Ana" }), "/api/auth/register", "POST", false, { nombre: "Ana" }],
    ["inicia sesión", () => login({ email: "ana@example.com" }), "/api/auth/login", "POST", false, { email: "ana@example.com" }],
    ["obtiene el usuario actual", () => obtenerUsuarioActual(), "/api/auth/me", undefined, true],
    ["obtiene gastos filtrados", () => obtenerGastos({ categoriaId: 3, desde: "", hasta: null, texto: "café" }), "/api/gastos?categoriaId=3&texto=caf%C3%A9", undefined, true],
    ["crea un gasto", () => crearGasto({ monto: 100 }), "/api/gastos", "POST", true, { monto: 100 }],
    ["actualiza un gasto", () => actualizarGasto(7, { monto: 100 }), "/api/gastos/7", "PUT", true, { monto: 100 }],
    ["elimina un gasto", () => eliminarGasto(7), "/api/gastos/7", "DELETE", true],
    ["obtiene categorías", () => obtenerCategorias(), "/api/categorias", undefined, true],
    ["crea una categoría", () => crearCategoria({ nombre: "Ocio" }), "/api/categorias", "POST", true, { nombre: "Ocio" }],
    ["actualiza una categoría", () => actualizarCategoria(4, { nombre: "Salud" }), "/api/categorias/4", "PUT", true, { nombre: "Salud" }],
    ["elimina una categoría", () => eliminarCategoria(4), "/api/categorias/4", "DELETE", true],
    ["obtiene el resumen filtrado", () => obtenerResumen({ desde: "2026-01-01", hasta: "2026-01-31" }), "/api/resumen?desde=2026-01-01&hasta=2026-01-31", undefined, true]
  ])("%s con el contrato HTTP esperado", async (_nombre, llamada, ruta, metodo, requiereAuth, cuerpoEsperado) => {
    guardarSesion({ token: "jwt-demo", user: { id: 1 } });
    fetch.mockResolvedValue(respuesta(metodo === "DELETE" ? null : { ok: true }, metodo === "DELETE" ? 204 : 200));

    await llamada();

    expect(fetch).toHaveBeenCalledTimes(1);
    const [url, opciones] = fetch.mock.calls[0];
    expect(url).toBe(ruta);
    expect(opciones.headers).toMatchObject({ "Content-Type": "application/json" });
    if (requiereAuth) expect(opciones.headers.Authorization).toBe("Bearer jwt-demo");
    else expect(opciones.headers.Authorization).toBeUndefined();
    expect(opciones.method ?? "GET").toBe(metodo ?? "GET");
    if (cuerpoEsperado !== undefined) {
      expect(JSON.parse(opciones.body)).toEqual(cuerpoEsperado);
    } else {
      expect(opciones.body).toBeUndefined();
    }
  });

  it("ante un 401 limpia la sesión, avisa a la aplicación y propaga el error", async () => {
    const alExpirarSesion = vi.fn();
    setUnauthorizedHandler(alExpirarSesion);
    guardarSesion({ token: "jwt-demo", user: { id: 1 } });
    fetch.mockResolvedValue(respuesta({ error: "Token inválido o expirado." }, 401));

    await expect(obtenerCategorias()).rejects.toThrow("Token inválido o expirado.");

    expect(localStorage.removeItem).toHaveBeenCalledTimes(2);
    expect(alExpirarSesion).toHaveBeenCalledTimes(1);
  });

  it("conserva el mensaje genérico cuando el servidor no entrega JSON", async () => {
    fetch.mockResolvedValue({ status: 500, ok: false, json: vi.fn().mockRejectedValue(new Error("sin JSON")) });

    await expect(obtenerCategorias()).rejects.toThrow("Ocurrió un error al comunicarse con el servidor.");
  });

  it("no agrega parámetros ni autorización cuando no se proporcionaron", async () => {
    fetch.mockResolvedValue(respuesta([]));

    await obtenerGastos({});

    expect(fetch).toHaveBeenCalledWith("/api/gastos", expect.objectContaining({
      headers: { "Content-Type": "application/json" }
    }));
  });
});
