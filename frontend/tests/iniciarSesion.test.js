import { describe, expect, it, vi } from "vitest";
import { iniciarSesion } from "../src/services/iniciarSesion";

describe("inicio de sesión", () => {
  it("autentica y guarda exactamente la sesión obtenida", async () => {
    const datos = { email: "ana@example.com", password: "12345678" };
    const sesion = { token: "jwt-demo", user: { id: 1, nombre: "Ana" } };
    const autenticar = vi.fn().mockResolvedValue(sesion);
    const guardarSesion = vi.fn();

    await expect(iniciarSesion(datos, { autenticar, guardarSesion })).resolves.toBe(sesion);

    expect(autenticar).toHaveBeenCalledTimes(1);
    expect(autenticar).toHaveBeenCalledWith(datos);
    expect(guardarSesion).toHaveBeenCalledTimes(1);
    expect(guardarSesion).toHaveBeenCalledWith(sesion);
  });

  it("propaga el error de autenticación sin guardar una sesión inexistente", async () => {
    const autenticar = vi.fn().mockRejectedValue(new Error("Credenciales inválidas"));
    const guardarSesion = vi.fn();

    await expect(iniciarSesion({}, { autenticar, guardarSesion })).rejects.toThrow("Credenciales inválidas");

    expect(guardarSesion).not.toHaveBeenCalled();
  });
});
