import { test, expect, listarGastos } from "./fixtures.js";
import { credenciales } from "./entorno.js";

test.beforeEach(async ({ page, datos }) => {
  await page.goto("/login");
  const { email, password } = credenciales();
  await page.getByLabel("Email", { exact: true }).fill(email);
  await page.getByLabel("Contraseña", { exact: true }).fill(password);
  await page.getByRole("button", { name: "Iniciar sesión", exact: true }).click();
  await expect(page).toHaveURL(/\/dashboard$/);
  await expect(page.getByRole("region", { name: "Nuevo gasto", exact: true })
    .getByRole("combobox", { name: "Categoría", exact: true }))
    .toContainText(datos.categoria.nombre);
});

function fila(page, descripcion) {
  return page.getByRole("row").filter({ has: page.getByRole("cell", { name: descripcion, exact: true }) });
}

async function crearDesdeFormulario(page, datos, descripcion) {
  const formulario = page.getByRole("region", { name: "Nuevo gasto", exact: true });
  await formulario.getByLabel("Descripción", { exact: true }).fill(descripcion);
  await formulario.getByLabel("Monto", { exact: true }).fill("123.45");
  await formulario.getByLabel("Fecha", { exact: true }).fill(datos.fecha);
  await formulario.getByRole("combobox", { name: "Categoría", exact: true }).selectOption(String(datos.categoria.id));
  await formulario.getByRole("button", { name: "Registrar gasto", exact: true }).click();
}

async function borrarDesdeFila(page, descripcion) {
  page.once("dialog", (dialog) => dialog.accept());
  await fila(page, descripcion).getByRole("button", { name: "Eliminar", exact: true }).click();
  await expect(fila(page, descripcion)).toHaveCount(0);
}

test("crear un gasto lo muestra en la lista y eliminarlo lo hace desaparecer", async ({ page, request, datos }) => {
  const descripcion = `${datos.prefijo} alta`;
  await crearDesdeFormulario(page, datos, descripcion);
  await expect(fila(page, descripcion)).toBeVisible();
  await expect(fila(page, descripcion).getByRole("cell", { name: /123,45/ })).toBeVisible();
  await borrarDesdeFila(page, descripcion);
  expect(await listarGastos(request, datos, datos.prefijo)).toEqual([]);
});

test("una descripción de espacios muestra el error y no crea ningún gasto", async ({ page, request, datos }) => {
  const antes = (await listarGastos(request, datos)).map((gasto) => gasto.id).sort((a, b) => a - b);
  // Tres espacios pasan minlength del HTML, pero la validación de la app los rechaza.
  await crearDesdeFormulario(page, datos, "   ");
  await expect(page.getByRole("alert")).toHaveText("La descripción debe tener entre 3 y 200 caracteres.");
  const despues = (await listarGastos(request, datos)).map((gasto) => gasto.id).sort((a, b) => a - b);
  expect(despues).toEqual(antes);
});

test("editar un gasto conserva los cambios al recargar la página", async ({ page, request, datos }) => {
  const descripcion = `${datos.prefijo} original`;
  const editado = `${datos.prefijo} editado`;
  await crearDesdeFormulario(page, datos, descripcion);
  await expect(fila(page, descripcion)).toBeVisible();
  await fila(page, descripcion).getByRole("button", { name: "Editar", exact: true }).click();
  const formulario = page.getByRole("region", { name: "Editar gasto", exact: true });
  await formulario.getByLabel("Descripción", { exact: true }).fill(editado);
  await formulario.getByLabel("Monto", { exact: true }).fill("87.65");
  await formulario.getByRole("button", { name: "Guardar cambios", exact: true }).click();
  await expect(fila(page, editado)).toBeVisible();
  await expect(fila(page, descripcion)).toHaveCount(0);
  await page.reload();
  await expect(fila(page, editado).getByRole("cell", { name: /87,65/ })).toBeVisible();
  await borrarDesdeFila(page, editado);
  expect(await listarGastos(request, datos, datos.prefijo)).toEqual([]);
});
