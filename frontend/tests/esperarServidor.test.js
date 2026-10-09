import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { esperarServidor } from "../src/services/esperarServidor";

const saludable = () => ({ ok: true, json: async () => ({ status: "healthy" }) });

describe("disponibilidad del servidor", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.stubGlobal("fetch", vi.fn());
  });
  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it("acepta el estado saludable real y no deja consultas periódicas activas", async () => {
    fetch.mockResolvedValue(saludable());
    await esperarServidor();
    expect(fetch).toHaveBeenCalledTimes(1);
    expect(fetch).toHaveBeenCalledWith("/health", expect.objectContaining({ method: "GET", cache: "no-store" }));
    expect(vi.getTimerCount()).toBe(0);
  });

  it.each([
    ["HTTP 503", { ok: false }],
    ["JSON que no confirma disponibilidad", { ok: true, json: async () => ({ status: "starting" }) }],
    ["HTML de arranque", { ok: true, json: async () => { throw new SyntaxError("no es JSON"); } }]
  ])("espera después de recibir %s y permite continuar cuando está saludable", async (_nombre, respuesta) => {
    fetch.mockResolvedValueOnce(respuesta).mockResolvedValue(saludable());
    const resultado = esperarServidor();
    await vi.advanceTimersByTimeAsync(3_000);
    await resultado;
    expect(fetch).toHaveBeenCalledTimes(2);
    expect(fetch.mock.calls.every(([url, options]) => url === "/health" && options.method === "GET")).toBe(true);
    expect(vi.getTimerCount()).toBe(0);
  });

  it("tolera un fallo transitorio de red sin reenviar operaciones de la app", async () => {
    fetch.mockRejectedValueOnce(new TypeError("sin conexión")).mockResolvedValue(saludable());
    const resultado = esperarServidor();
    await vi.advanceTimersByTimeAsync(3_000);
    await resultado;
    expect(fetch.mock.calls.map(([url]) => url)).toEqual(["/health", "/health"]);
  });

  it("cancela una consulta colgada y vuelve a comprobar el servicio", async () => {
    fetch.mockImplementationOnce((_url, { signal }) => new Promise((_resolve, reject) => {
      signal.addEventListener("abort", () => reject(signal.reason), { once: true });
    })).mockResolvedValue(saludable());
    const resultado = esperarServidor();
    await vi.advanceTimersByTimeAsync(18_000);
    await resultado;
    expect(fetch).toHaveBeenCalledTimes(2);
    expect(fetch.mock.calls[0][1].signal.aborted).toBe(true);
    expect(vi.getTimerCount()).toBe(0);
  });

  it("deja de consultar al superar el límite total y comunica que se puede reintentar", async () => {
    fetch.mockResolvedValue({ ok: false });
    const resultado = esperarServidor();
    const comprobacion = expect(resultado).rejects.toThrow("El servidor no respondió a tiempo.");
    await vi.advanceTimersByTimeAsync(120_000);
    await comprobacion;
    const consultas = fetch.mock.calls.length;
    await vi.advanceTimersByTimeAsync(30_000);
    expect(fetch).toHaveBeenCalledTimes(consultas);
    expect(vi.getTimerCount()).toBe(0);
  });

  it("cancela la espera entre consultas cuando se abandona la pantalla", async () => {
    const controller = new AbortController();
    fetch.mockResolvedValue({ ok: false });
    const resultado = esperarServidor({ signal: controller.signal });
    const comprobacion = expect(resultado).rejects.toMatchObject({ name: "AbortError" });
    await vi.advanceTimersByTimeAsync(0);
    controller.abort();
    await comprobacion;
    expect(fetch).toHaveBeenCalledTimes(1);
    expect(vi.getTimerCount()).toBe(0);
  });

  it("no hace consultas si la operación ya fue cancelada", async () => {
    const controller = new AbortController();
    controller.abort();
    await expect(esperarServidor({ signal: controller.signal })).rejects.toMatchObject({ name: "AbortError" });
    expect(fetch).not.toHaveBeenCalled();
    expect(vi.getTimerCount()).toBe(0);
  });
});
