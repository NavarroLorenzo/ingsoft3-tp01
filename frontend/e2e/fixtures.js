import { randomUUID } from "node:crypto";
import { test as base, expect } from "@playwright/test";
import { API, credenciales } from "./entorno.js";

export { expect };

export async function listarGastos(request, datos, texto) {
  const response = await request.get(`${API}/api/gastos`, {
    headers: datos.headers,
    params: texto === undefined ? {} : { texto }
  });
  expect(response.status()).toBe(200);
  return response.json();
}

export const test = base.extend({
  datos: async ({ request }, use) => {
    const login = await request.post(`${API}/api/auth/login`, { data: credenciales() });
    expect(login.status(), "La cuenta de pruebas debe estar preparada antes de ejecutar los specs.").toBe(200);
    const { token } = await login.json();
    const headers = { Authorization: `Bearer ${token}` };
    const categorias = await request.get(`${API}/api/categorias`, { headers });
    expect(categorias.status()).toBe(200);
    const categoria = (await categorias.json()).find((item) => item.nombre === "Comida");
    expect(categoria, "La categoría inicial Comida debe existir.").toBeDefined();
    const datos = {
      headers,
      categoria,
      prefijo: `tp7-${Date.now()}-${randomUUID().slice(0, 8)}`,
      fecha: new Date().toISOString().slice(0, 10)
    };
    try {
      await use(datos);
    } finally {
      // Busca sólo los datos propios: también limpia si un test falla después del alta.
      for (const gasto of await listarGastos(request, datos, datos.prefijo)) {
        const response = await request.delete(`${API}/api/gastos/${gasto.id}`, { headers });
        expect(response.status(), "La limpieza debe borrar el gasto propio.").toBe(204);
      }
      expect(await listarGastos(request, datos, datos.prefijo)).toEqual([]);
    }
  }
});
