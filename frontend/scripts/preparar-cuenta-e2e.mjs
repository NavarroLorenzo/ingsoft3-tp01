import { API, credenciales } from "../e2e/entorno.js";

// Preparación explícita: se reutiliza la cuenta; los specs no crean usuarios por corrida.
const datos = credenciales();
const enviar = (ruta, body) => fetch(`${API}${ruta}`, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(body),
  signal: AbortSignal.timeout(60_000)
});
const login = await enviar("/api/auth/login", datos);
if (login.status === 200) {
  console.log("La cuenta de pruebas ya existe y permite iniciar sesión.");
} else if (login.status === 401) {
  const registro = await enviar("/api/auth/register", { nombre: "Pruebas TP7", ...datos });
  if (registro.status !== 201) {
    throw new Error(`No se pudo preparar la cuenta (HTTP ${registro.status}). Verificá las credenciales.`);
  }
  console.log("Cuenta de pruebas creada. Se reutiliza en las siguientes corridas.");
} else {
  throw new Error(`La API no permitió iniciar sesión (HTTP ${login.status}).`);
}
