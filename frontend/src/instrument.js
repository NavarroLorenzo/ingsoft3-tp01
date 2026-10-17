import * as Sentry from "@sentry/react";

// Se importa antes de React para registrar errores y tiempos desde que abre la app.
Sentry.init({
  // Este DSN es público: permite enviar eventos al proyecto del frontend, no leerlos.
  dsn: "https://6263155857d38f405d744caa93f12e54@o4512273864196096.ingest.us.sentry.io/4512273878745088",
  environment: window.location.hostname,
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0, // La guía pide medir todos los pedidos durante el ejercicio.
  // Para medir rutas y errores no necesitamos credenciales ni el contenido de los gastos.
  dataCollection: {
    userInfo: false,
    cookies: false,
    httpHeaders: false,
    httpBodies: [],
    urlQueryParams: false,
  },
});
