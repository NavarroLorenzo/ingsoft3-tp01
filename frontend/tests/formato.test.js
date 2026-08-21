import { describe, expect, it } from "vitest";
import { formatearFecha, formatearMoneda, validarEmail, validarGasto, validarRegistro } from "../src/utils/formato";

describe("formato", () => {
  it("formatea moneda argentina", () => {
    expect(formatearMoneda(25400.5)).toContain("25.400,50");
  });

  it("formatea fechas como DD/MM/YYYY", () => {
    expect(formatearFecha("2026-08-12")).toBe("12/08/2026");
  });

  it("valida un gasto simple", () => {
    expect(validarGasto({ descripcion: "ab", monto: 100, fecha: "2026-08-12", categoriaId: 1 })).toContain("descripción");
    expect(validarGasto({ descripcion: "Supermercado", monto: 100.5, fecha: "2026-08-12", categoriaId: 1 })).toBeNull();
  });

  it("valida email y datos de registro", () => {
    expect(validarEmail("correo@ejemplo.com")).toBe(true);
    expect(validarEmail("correo-invalido")).toBe(false);
    expect(validarRegistro({ nombre: "Ana", email: "ana@ejemplo.com", password: "12345678", confirmacion: "12345678" })).toBeNull();
    expect(validarRegistro({ nombre: "Ana", email: "ana@ejemplo.com", password: "12345678", confirmacion: "otra-clave" })).toContain("coinciden");
  });
});
