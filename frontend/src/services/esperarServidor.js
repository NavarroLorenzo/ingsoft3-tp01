function pausa(ms, signal) {
  signal.throwIfAborted();
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      signal.removeEventListener("abort", cancelar);
      resolve();
    }, ms);
    function cancelar() {
      clearTimeout(timer);
      signal.removeEventListener("abort", cancelar);
      reject(signal.reason);
    }
    signal.addEventListener("abort", cancelar, { once: true });
  });
}

// Sólo se reintenta el GET de salud; login, registro y escrituras se envían una vez.
export async function esperarServidor({ signal } = {}) {
  const limite = new AbortController();
  const espera = AbortSignal.any([limite.signal, ...(signal ? [signal] : [])]);
  const timerLimite = setTimeout(() => {
    limite.abort(new Error("El servidor no respondió a tiempo. Volvé a intentar en unos minutos."));
  }, 120_000);

  try {
    espera.throwIfAborted();
    // Una sola petición directa, sin cookies, para activar la API en Render.
    // La redirección sale de Nginx según el entorno; no se interpreta su respuesta opaca.
    void fetch("/health-wakeup", {
      method: "GET",
      mode: "no-cors",
      credentials: "omit",
      referrerPolicy: "no-referrer",
      cache: "no-store",
      signal: espera
    }).catch(() => {});

    while (true) {
      espera.throwIfAborted();
      const intento = new AbortController();
      const timerIntento = setTimeout(() => intento.abort(), 15_000);
      try {
        const response = await fetch("/health", {
          method: "GET",
          cache: "no-store",
          signal: AbortSignal.any([espera, intento.signal])
        });
        const body = response.ok ? await response.json() : null;
        espera.throwIfAborted();
        if (body?.status === "healthy") return;
      } catch {
        // Una respuesta HTML, un 503 o un fallo de red todavía no indican disponibilidad.
        espera.throwIfAborted();
      } finally {
        clearTimeout(timerIntento);
      }
      await pausa(3_000, espera);
    }
  } finally {
    clearTimeout(timerLimite);
    limite.abort(); // Cancela también el pedido de arranque si aún sigue pendiente.
  }
}
