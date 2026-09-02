3.7 Publicar las imágenes en un registry
Riel canónico — GitHub Container Registry (ghcr.io), gratis para imágenes públicas.

Paso 1 — generá el token. Para publicar en ghcr necesitás un Personal Access Token clásico con el permiso write:packages. Se crea así, una sola vez:

GitHub → tu foto (arriba a la derecha) → Settings → abajo de todo, Developer settings → Personal access tokens → Tokens (classic) → Generate new token (classic). Ponele un nombre (TP2 ghcr), una expiración, y tildá write:packages (al tildarlo se marca solo read:packages y repo). Generate token → copiá el token ahora, porque no se vuelve a mostrar.

🔴 Tiene que ser classic, no fine-grained: los fine-grained no funcionan con ghcr — el docker login te dice Succeeded y después el push falla con denied: permission_denied. Es el error más confuso de este paso.

Paso 2 — logueate al registry. Docker te va a pedir una contraseña: pegá el token, no tu contraseña de GitHub (no la vas a ver mientras la pegás, es normal).

docker login ghcr.io -u <tu_usuario>
# Password: ← pegás el token y Enter
# → Login Succeeded
💡 ¿Preferís no tipearlo cada vez? Guardalo en una variable y pasáselo por entrada estándar — así el token no queda en el historial del shell:

export CR_PAT=ghp_xxxxxxxx
echo $CR_PAT | docker login ghcr.io -u <tu_usuario> --password-stdin
Si tenés el GitHub CLI y preferís usar su token: gh auth refresh -h github.com -s write:packages y después gh auth token | docker login ghcr.io -u <tu_usuario> --password-stdin. Ojo que el refresh no siempre toma el permiso nuevo; si el push sigue dando denied, volvé al token clásico de arriba.

Paso 3 — tag + push. ⚠️ <tu_usuario> va todo en MINÚSCULAS, aunque tu usuario de GitHub tenga mayúsculas: Docker no las acepta en el nombre de una imagen y corta con invalid reference format: repository name must be lowercase. (El docker login sí las acepta, lo que hace el error todavía más desconcertante.)

docker tag mi-backend:dev ghcr.io/<tu_usuario>/mi-backend:v0.1.0
docker tag mi-frontend:dev ghcr.io/<tu_usuario>/mi-frontend:v0.1.0

docker images | grep -E 'ghcr|mi-backend|mi-frontend'    # (PowerShell: | Select-String 'ghcr|mi-')

docker push ghcr.io/<tu_usuario>/mi-backend:v0.1.0
docker push ghcr.io/<tu_usuario>/mi-frontend:v0.1.0
💡 Mirá el identificador de las imágenes en ese último listado y compará con mi-backend:dev: es el mismo. docker tag no copia nada — le agrega un nombre más a una imagen que ya existe. Ese detalle vuelve en el Paso 5.

Paso 4 — hacelas públicas. 🔴 Los packages de ghcr nacen privados, y mientras lo estén nadie puede hacer docker pull de tu imagen: ni la cátedra, ni otra máquina, ni un pipeline, ni un compose que la referencie. Es el tropiezo más común de esta sección, y hay que hacerlo para las dos imágenes:

Tu perfil de GitHub → pestaña Packages → clic en el package → Package settings (abajo a la derecha) → Change visibility → Public → confirmar escribiendo el nombre.

💡 La imagen recién pusheada no aparece dentro del repositorio, así que buscala en tu perfil. Si querés que quede linkeada al repo (útil para el TP7), agregale al Dockerfile la línea LABEL org.opencontainers.image.source=https://github.com/<tu_usuario>/<tu_repo>.

Paso 5 — comprobá que quedó público de verdad. Decir que está público es fácil; la prueba es bajarla sin credenciales:

docker compose down                                     # si tenés algo corriendo con esa imagen
docker logout ghcr.io                                   # dejás de estar autenticado
docker rmi ghcr.io/<tu_usuario>/mi-backend:v0.1.0 mi-backend:dev
docker pull ghcr.io/<tu_usuario>/mi-backend:v0.1.0      # …y la bajás de cero
🔴 Fijate que el rmi lleva LOS DOS nombres. Como docker tag no copió nada, esa imagen tiene dos nombres y un solo cuerpo: si borrás sólo el del registro, Docker te contesta Untagged y no borra nada.

⚠️ Y aun así el pull va a decir Already exists en todas las capas. Está bien, no fallaste. Docker no guarda imágenes: guarda capas, identificadas por su contenido. Tu docker compose up -d --build de §3.6 construyó su propia imagen del backend (se llama <tu-carpeta>-backend) a partir del mismo código, así que sus capas son las mismas y siguen en tu disco mientras esa imagen exista. Borrar los dos nombres no borra el cuerpo si alguien más lo está usando.

Lo que este paso prueba —y es lo que importa— es que pudiste pedir la imagen sin credenciales, estando deslogueado. La descarga de verdad la vas a ver en el Paso 6, donde también se van las imágenes del compose.

Si el pull funciona sin sesión, cualquiera puede correr tu imagen — eso es el checkpoint, no que la página diga Public.

⚠️ Una advertencia honesta sobre arquitecturas. La imagen que publicaste sirve para máquinas con el mismo tipo de procesador que la tuya: si la construiste en una Mac moderna es para ARM, y si la construiste en una PC común es para Intel/AMD. Alguien con la otra arquitectura recibe no matching manifest for linux/amd64 in the manifest list entries — y los runners de CI del TP7 son Intel. Para este TP alcanza con saberlo (anotalo en decisiones.md y decí en qué máquina la construiste); en el TP7 lo vamos a resolver con docker buildx, que construye para las dos a la vez.

Paso 6 — ¿y para qué publicaste? Para que el sistema se pueda levantar sin tu código. El docker-compose.yml que escribiste usa build: ./backend, así que necesita el repositorio. Escribí una variante que use image: con el nombre completo, y probala:

docker-compose.registry.yml, en la raíz, al lado del otro. Es el mismo compose: mismos servicios, mismo volumen, mismo healthcheck, mismo depends_on. Lo único que cambia son las dos líneas donde decía build::

services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: app
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      retries: 10

  backend:
    image: ghcr.io/<tu_usuario>/mi-backend:v0.1.0     # ← antes: build: ./backend
    environment:
      ConnectionStrings__Default: "Host=db;Database=app;Username=postgres;Password=${DB_PASSWORD}"
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy

  frontend:
    image: ghcr.io/<tu_usuario>/mi-frontend:v0.1.0    # ← antes: build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

volumes:
  db_data:
bat --paging=never docker-compose.registry.yml   # compará con el otro: cambian dos líneas
docker compose down --rmi local     # ← el flag se lleva TAMBIÉN las imágenes que compose construyó
docker rmi ghcr.io/<tu_usuario>/mi-backend:v0.1.0 ghcr.io/<tu_usuario>/mi-frontend:v0.1.0 mi-frontend:dev
docker builder prune -af            # ← y esto vacía el CACHE DE CONSTRUCCIÓN, que también las guarda
docker compose -f docker-compose.registry.yml up -d    # NO construye: descarga, y AHORA sí se ve
🔴 Las capas se esconden en TRES lugares, y hay que sacarlas de los tres. Docker no guarda imágenes: guarda capas, identificadas por su contenido, y sólo las libera cuando ya nadie las referencia.

La imagen que construyó el compose de §3.6 (<tu-carpeta>-backend) → la saca el --rmi local.
Los nombres que le pusiste vos (mi-backend:dev, el del registro) → los saca el rmi.
El cache de construcción → lo saca el builder prune. Éste es el que nadie ve venir: puede tener más de un giga de capas guardadas para acelerar tus próximos builds.
Si te salteás cualquiera de los tres, el up te contesta Already exists en todas las capas y no baja nada — el mismo malentendido del Paso 5. Con los tres, la descarga se ve capa por capa.

⚠️ Después del builder prune, tu próximo docker build va a tardar como el primero: se quedó sin cache. Es el precio de hacer la prueba en serio, y la hacés una sola vez. (mi-backend:dev no está en la lista del rmi porque ya lo borraste en el Paso 5; si lo ponés, Docker te contesta No such image.)

⚠️ Este compose también necesita el .env. Usa ${DB_PASSWORD} igual que el otro, así que si se lo pasás a alguien "para levantar sin el código", tiene que ir acompañado de la plantilla —o de la contraseña por otro canal—. Son dos archivos, no uno. Eso no es un defecto: es la disciplina del secreto que no viaja por el repositorio.

El sample tiene este archivo resuelto en su rama main, para que compares. Es el punto de partida del TP7: allá el pipeline construye y publica, y el entorno consume exactamente así.

⚠️ Un gotcha más de ghcr: el docker login da OK con cualquier token válido, tenga o no el permiso de packages… y recién falla el push, con denied — el login exitoso no garantiza el permiso.

¿Por qué ghcr y no Docker Hub? Los dos sirven y aceptamos cualquiera de los dos. Elegimos ghcr para la cátedra porque: (1) ya tenés la cuenta —es la de GitHub del TP1—; (2) las imágenes quedan junto al código, así tu entrega es un solo lugar; y (3) en el TP7, cuando publique el pipeline en vez de vos, GitHub Actions puede autenticarse contra ghcr sin secretos: usa el GITHUB_TOKEN del propio workflow, declarándole permissions: packages: write. Docker Hub también se automatiza sin problema — solo que ahí hay que crear un token y guardarlo como secreto del repo.

Si preferís Docker Hub: docker login (sin servidor, es el default) + docker tag mi-backend:dev <usuario>/mi-backend:v0.1.0 + docker push <usuario>/mi-backend:v0.1.0. Ventaja: los repos públicos nacen públicos, así que te ahorrás el paso de cambiar visibilidad. Desventajas menores: las cuentas gratuitas tienen límite de descargas por hora (afecta al pull, no al push) y solo permiten un repositorio privado.

✅ Checkpoint: las imágenes son visibles en tu perfil (pestaña Packages en GitHub), con visibilidad pública, y cualquiera puede hacer docker pull de ellas.