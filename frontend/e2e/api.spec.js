import { test, expect, listarGastos } from "./fixtures.js";
import { API } from "./entorno.js";

function gastoValido(datos) {
  return { descripcion: `${datos.prefijo} gasto`, monto: 123.45, fecha: datos.fecha, categoriaId: datos.categoria.id };
}

test("el alta persiste un gasto y el borrado lo quita de la base", async ({ request, datos }) => {
  const input = gastoValido(datos);
  const alta = await request.post(`${API}/api/gastos`, { headers: datos.headers, data: input });
  expect(alta.status()).toBe(201);
  const creado = await alta.json();
  expect(creado.id).toBeGreaterThan(0);
  expect(await listarGastos(request, datos, datos.prefijo)).toEqual([
    expect.objectContaining({ id: creado.id, ...input })
  ]);
  const borrado = await request.delete(`${API}/api/gastos/${creado.id}`, { headers: datos.headers });
  expect(borrado.status()).toBe(204);
  expect(await listarGastos(request, datos, datos.prefijo)).toEqual([]);
});

test("una descripción vacía devuelve 400 y no inserta un gasto", async ({ request, datos }) => {
  const antes = (await listarGastos(request, datos)).map((gasto) => gasto.id).sort((a, b) => a - b);
  const alta = await request.post(`${API}/api/gastos`, {
    headers: datos.headers,
    data: { ...gastoValido(datos), descripcion: "" }
  });
  expect(alta.status()).toBe(400);
  expect(await alta.json()).toEqual({ error: "La descripción es obligatoria." });
  const despues = (await listarGastos(request, datos)).map((gasto) => gasto.id).sort((a, b) => a - b);
  expect(despues).toEqual(antes);
});

test("editar un gasto persiste el monto y la descripción nuevos", async ({ request, datos }) => {
  const input = gastoValido(datos);
  const alta = await request.post(`${API}/api/gastos`, { headers: datos.headers, data: input });
  expect(alta.status()).toBe(201);
  const { id } = await alta.json();
  const editado = { ...input, descripcion: `${datos.prefijo} editado`, monto: 87.65 };
  const cambio = await request.put(`${API}/api/gastos/${id}`, { headers: datos.headers, data: editado });
  expect(cambio.status()).toBe(200);
  const consulta = await request.get(`${API}/api/gastos/${id}`, { headers: datos.headers });
  expect(consulta.status()).toBe(200);
  expect(await consulta.json()).toEqual(expect.objectContaining({ id, ...editado }));
  const borrado = await request.delete(`${API}/api/gastos/${id}`, { headers: datos.headers });
  expect(borrado.status()).toBe(204);
  expect(await listarGastos(request, datos, datos.prefijo)).toEqual([]);
});
