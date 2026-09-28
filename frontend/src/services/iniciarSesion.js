// La API y la persistencia entran como dependencias para que esta regla se
// pueda probar sin red ni localStorage reales.
export async function iniciarSesion(datos, { autenticar, guardarSesion }) {
  const sesion = await autenticar(datos);
  guardarSesion(sesion);
  return sesion;
}
