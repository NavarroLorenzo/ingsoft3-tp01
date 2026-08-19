import { beforeEach, describe, expect, it, vi } from "vitest";
import { getAuthHeaders, limpiarSesion, guardarSesion, setUnauthorizedHandler } from "../src/api/api";

describe("sesión de autenticación", () => {
  beforeEach(() => {
    const store = {};
    global.localStorage = { getItem: vi.fn((key) => store[key] || null), setItem: vi.fn((key, value) => { store[key] = value; }), removeItem: vi.fn((key) => { delete store[key]; }) };
    setUnauthorizedHandler(null);
  });

  it("agrega Authorization y limpia la sesión", () => {
    guardarSesion({ token: "jwt-demo", user: { id: 1, nombre: "Ana" } });
    expect(getAuthHeaders()).toEqual({ Authorization: "Bearer jwt-demo" });
    limpiarSesion();
    expect(getAuthHeaders()).toEqual({});
  });
});
