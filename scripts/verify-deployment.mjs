import { appendFile } from "node:fs/promises";
import { pathToFileURL } from "node:url";

// El mismo control se usa en QA y PROD. No necesita claves ni dependencias npm.
export async function verifyDeployment({
  apiUrl, frontendUrl, expectedCommit,
  attempts = 30, intervalMs = 20_000, requestTimeoutMs = 10_000,
  fetchImpl = globalThis.fetch,
  pause = (ms) => new Promise((resolve) => setTimeout(resolve, ms)),
  log = console.log,
}) {
  if (!/^[0-9a-f]{40}$/.test(expectedCommit || "")) {
    throw new Error("EXPECTED_COMMIT debe ser el SHA completo de 40 caracteres.");
  }
  for (const base of [apiUrl, frontendUrl]) {
    const url = new URL(base);
    if (!["http:", "https:"].includes(url.protocol) || url.username || url.password) {
      throw new Error("Las URLs de despliegue deben ser HTTP(S) y no contener credenciales.");
    }
  }
  const api = apiUrl.replace(/\/+$/, "");
  const front = frontendUrl.replace(/\/+$/, "");
  const options = () => ({
    cache: "no-store", redirect: "error",
    headers: { "Cache-Control": "no-cache" },
    signal: AbortSignal.timeout(requestTimeoutMs),
  });
  async function readVersion(url) {
    const response = await fetchImpl(`${url}?deploy=${expectedCommit}`, options());
    if (response.status !== 200) throw new Error(`${new URL(url).pathname}: HTTP ${response.status}`);
    try { return await response.json(); }
    catch { throw new Error(`${new URL(url).pathname} no devolvió JSON de versión.`); }
  }
  async function status(url) {
    const response = await fetchImpl(url, options());
    await response.body?.cancel();
    return response.status;
  }

  log(`Esperando backend y frontend del commit ${expectedCommit}...`);
  let lastError = "sin respuesta";
  for (let attempt = 1; attempt <= attempts; attempt++) {
    try {
      const [backend, frontend] = await Promise.all([
        readVersion(`${api}/health`), readVersion(`${front}/version.json`),
      ]);
      if (backend?.status !== "healthy" || backend?.commit !== expectedCommit || frontend?.commit !== expectedCommit) {
        throw new Error(`backend=${backend?.commit ?? "sin commit"}, frontend=${frontend?.commit ?? "sin commit"}; deben coincidir con ${expectedCommit}.`);
      }
      const [home, proxy] = await Promise.all([status(`${front}/`), status(`${front}/api/categorias`)]);
      if (home !== 200 || proxy !== 401) {
        throw new Error(`Frontend HTTP ${home}, proxy HTTP ${proxy}; se esperaban 200 y 401.`);
      }
      log(`Versión verificada en ambos servicios: ${expectedCommit}. Base de datos y proxy disponibles.`);
      return { backendCommit: backend.commit, frontendCommit: frontend.commit };
    } catch (error) {
      lastError = error.message;
      log(`Intento ${attempt}/${attempts}: ${lastError}`);
    }
    if (attempt < attempts) await pause(intervalMs);
  }
  throw new Error(`No se confirmó el despliegue de ${expectedCommit}. ${lastError}`);
}

async function main() {
  const result = await verifyDeployment({
    apiUrl: process.env.URL_API,
    frontendUrl: process.env.URL_FRONT,
    expectedCommit: process.env.EXPECTED_COMMIT,
  });
  if (process.env.GITHUB_STEP_SUMMARY) {
    await appendFile(process.env.GITHUB_STEP_SUMMARY, [
      `## Versión desplegada verificada (${process.env.DEPLOY_ENV || "entorno"})`, "",
      "| Servicio | Commit observado |", "|---|---|",
      `| Backend: /health | \`${result.backendCommit}\` |`,
      `| Frontend: /version.json | \`${result.frontendCommit}\` |`,
      "", "Los dos coinciden con el commit de la corrida; PostgreSQL y el proxy responden.", "",
    ].join("\n"));
  }
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  main().catch((error) => { console.error(error.message); process.exitCode = 1; });
}
