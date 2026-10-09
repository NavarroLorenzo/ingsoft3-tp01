import { existsSync } from "node:fs";
import { fileURLToPath } from "node:url";

const archivo = fileURLToPath(new URL("../.env.e2e.local", import.meta.url));
if (existsSync(archivo)) process.loadEnvFile(archivo);

export const API = (process.env.API_BASE_URL || "http://localhost:8080").replace(/\/$/, "");
export const FRONT = process.env.E2E_BASE_URL || "http://localhost:3000";

export function credenciales() {
  const email = process.env.E2E_EMAIL;
  const password = process.env.E2E_PASSWORD;
  if (!email || !password) {
    throw new Error("Definí E2E_EMAIL y E2E_PASSWORD para una cuenta exclusiva de pruebas.");
  }
  return { email, password };
}
