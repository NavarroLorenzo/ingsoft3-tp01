import { afterEach, describe, expect, it, vi } from "vitest";
import { fechaActual, formatearFecha, formatearMoneda, validarCategoria, validarEmail, validarFechaGasto, validarGasto, validarRegistro } from "../src/utils/formato";

const gastoValido = {
  descripcion: "Supermercado",
  monto: 100.5,
  fecha: "2026-08-12",
  categoriaId: 1
};

describe("formato", () => {
  afterEach(() => vi.useRealTimers());

  it("formatea moneda argentina", () => {
    expect(formatearMoneda(25400.5)).toContain("25.400,50");
  });

  it("formatea fechas como DD/MM/YYYY", () => {
    expect(formatearFecha("2026-08-12")).toBe("12/08/2026");
    expect(formatearFecha("")).toBe("");
  });

  it("obtiene la fecha actual sin depender del reloj del equipo", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-08-12T12:00:00Z"));

    expect(fechaActual()).toBe("2026-08-12");
  });

  it.each([
    ["fecha no provista", undefined, "obligatoria"],
    ["fecha vacía", "", "obligatoria"],
    ["formato distinto de ISO", "12/08/2026", "formato"],
    ["día inexistente", "2026-02-30", "fecha válida"]
  ])("rechaza una %s", (_nombre, fecha, mensaje) => {
    expect(validarFechaGasto(fecha)).toContain(mensaje);
  });

  it("acepta una fecha ISO real, incluso en año bisiesto", () => {
    expect(validarFechaGasto("2028-02-29")).toBeNull();
  });

  it.each([
    ["descripción demasiado corta", { descripcion: "ab" }, "descripción"],
    ["monto cero", { monto: 0 }, "monto"],
    ["monto con tres decimales", { monto: 1.123 }, "decimales"],
    ["fecha ausente", { fecha: "" }, "fecha"],
    ["categoría ausente", { categoriaId: "" }, "categoría"]
  ])("rechaza un gasto con %s", (_nombre, cambios, mensaje) => {
    expect(validarGasto({ ...gastoValido, ...cambios })).toContain(mensaje);
  });

  it("acepta un gasto válido y quita los espacios de la descripción", () => {
    expect(validarGasto({ ...gastoValido, descripcion: "  Supermercado  " })).toBeNull();
  });

  it.each([
    ["una letra", "A", true],
    ["dos letras", "AB", false],
    ["cincuenta letras", "a".repeat(50), false],
    ["cincuenta y una letras", "a".repeat(51), true]
  ])("valida una categoría de %s", (_nombre, valor, invalida) => {
    const resultado = validarCategoria(`  ${valor}  `);
    if (invalida) expect(resultado).toContain("nombre");
    else expect(resultado).toBeNull();
  });

  it("valida email y datos de registro", () => {
    expect(validarEmail("correo@ejemplo.com")).toBe(true);
    expect(validarEmail("correo-invalido")).toBe(false);
    expect(validarRegistro({ nombre: "Ana", email: "ana@ejemplo.com", password: "12345678", confirmacion: "12345678" })).toBeNull();
    expect(validarRegistro({ nombre: "A", email: "ana@ejemplo.com", password: "12345678", confirmacion: "12345678" })).toContain("nombre");
    expect(validarRegistro({ nombre: "Ana", email: "correo-invalido", password: "12345678", confirmacion: "12345678" })).toContain("email");
    expect(validarRegistro({ nombre: "Ana", email: "ana@ejemplo.com", password: "corta", confirmacion: "corta" })).toContain("contraseña");
    expect(validarRegistro({ nombre: "Ana", email: "ana@ejemplo.com", password: "12345678", confirmacion: "otra-clave" })).toContain("coinciden");
  });
});
