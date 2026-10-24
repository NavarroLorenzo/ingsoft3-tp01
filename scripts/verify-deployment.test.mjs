import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { test } from "node:test";
import { fileURLToPath } from "node:url";
import { verifyDeployment } from "./verify-deployment.mjs";

const COMMIT = "0123456789abcdef0123456789abcdef01234567";
const OLD = "abcdef0123456789abcdef0123456789abcdef01";
const config = {
  apiUrl: "https://api.example.test", frontendUrl: "https://front.example.test",
  expectedCommit: COMMIT, attempts: 1, log() {}, pause: async () => {},
};

function service({ backendCommit = COMMIT, frontendCommit = COMMIT, database = "healthy", apiStatus = 200, proxy = 401, htmlVersion = false } = {}) {
  return async (url) => {
    const { pathname } = new URL(url);
    if (pathname === "/health") return Response.json({ status: database, commit: backendCommit }, { status: apiStatus });
    if (pathname === "/version.json") return htmlVersion ? new Response("<html>versión anterior sin metadata</html>") : Response.json({ commit: frontendCommit });
    if (pathname === "/api/categorias") return new Response(null, { status: proxy });
    return new Response("<html>app</html>");
  };
}

test("acepta únicamente las dos versiones esperadas, con salud y proxy correctos", async () => {
  assert.deepEqual(await verifyDeployment({ ...config, fetchImpl: service() }), {
    backendCommit: COMMIT, frontendCommit: COMMIT,
  });
});

test("un backend anterior saludable no hace pasar el deploy", async () => {
  await assert.rejects(verifyDeployment({ ...config, fetchImpl: service({ backendCommit: OLD }) }), /No se confirmó/);
});

test("el backend nuevo no alcanza si el frontend sigue en otra versión", async () => {
  await assert.rejects(verifyDeployment({ ...config, fetchImpl: service({ frontendCommit: OLD }) }), /No se confirmó/);
});

test("la SPA devolviendo HTML 200 en version.json no se acepta como versión", async () => {
  await assert.rejects(verifyDeployment({ ...config, fetchImpl: service({ htmlVersion: true }) }), /no devolvió JSON/);
});

test("no acepta una base caída aunque el commit sea el correcto", async () => {
  await assert.rejects(verifyDeployment({ ...config, fetchImpl: service({ database: "unhealthy", apiStatus: 503 }) }), /HTTP 503/);
});

test("conserva la comprobación del proxy protegido", async () => {
  await assert.rejects(verifyDeployment({ ...config, fetchImpl: service({ proxy: 502 }) }), /proxy HTTP 502/);
});

test("reintenta la versión anterior y sólo termina cuando aparecen ambas nuevas", async () => {
  let attempt = 0;
  let pauses = 0;
  const fetchImpl = async (url, options) => {
    if (new URL(url).pathname === "/health") attempt++;
    return service({ backendCommit: attempt === 1 ? OLD : COMMIT })(url, options);
  };
  const result = await verifyDeployment({ ...config, attempts: 3, fetchImpl, pause: async () => { pauses++; } });
  assert.equal(result.backendCommit, COMMIT);
  assert.equal(attempt, 2);
  assert.equal(pauses, 1);
});

test("agota los reintentos sin convertir la falta de versión en éxito", async () => {
  let calls = 0;
  await assert.rejects(verifyDeployment({ ...config, attempts: 2, fetchImpl: async () => {
    calls++;
    throw new Error("servicio iniciando");
  } }), /No se confirmó/);
  assert.equal(calls, 4);
});

test("rechaza un SHA abreviado antes de hacer requests", async () => {
  let called = false;
  await assert.rejects(verifyDeployment({ ...config, expectedCommit: "0123456", fetchImpl: async () => { called = true; } }), /40 caracteres/);
  assert.equal(called, false);
});

test("el comando devuelve salida de error cuando falta una configuración válida", () => {
  const result = spawnSync(process.execPath, [fileURLToPath(new URL("./verify-deployment.mjs", import.meta.url))], {
    encoding: "utf8", env: { ...process.env, EXPECTED_COMMIT: "incorrecto" },
  });
  assert.equal(result.status, 1);
  assert.match(result.stderr, /40 caracteres/);
});
