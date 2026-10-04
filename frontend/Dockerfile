# ---- Etapa 1: build (imagen con Node.js) ----
FROM node:22-alpine AS build

WORKDIR /app

# Se copian primero los manifiestos de dependencias para aprovechar la caché.
COPY package.json package-lock.json ./

# Instala exactamente las dependencias definidas en package-lock.json.
RUN npm ci

# Copia el código fuente una vez resueltas las dependencias.
COPY . .

# Genera los archivos estáticos optimizados para producción.
RUN npm run build

# ---- Etapa 2: tests y cobertura (reutiliza Node y node_modules) ----
FROM build AS test

ENTRYPOINT ["npm", "run", "test:coverage"]

# ---- Etapa 3: runtime (servidor web liviano, sin Node.js) ----
FROM nginx:1.29-alpine

# Configura Nginx al arrancar con la URL del backend de cada entorno.
COPY default.conf.template /etc/nginx/templates/default.conf.template

# Copia únicamente los archivos estáticos generados en la etapa de build.
COPY --from=build /app/dist /usr/share/nginx/html

# Docker Compose usa estos valores; Render los reemplaza para QA y PROD.
ENV BACKEND_URL=http://backend:8080 \
    DNS_RESOLVER=127.0.0.11

# Documenta el puerto HTTP expuesto por Nginx.
EXPOSE 80
