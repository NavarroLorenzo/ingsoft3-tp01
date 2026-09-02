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

# ---- Etapa 2: runtime (servidor web liviano, sin Node.js) ----
FROM nginx:1.29-alpine

LABEL org.opencontainers.image.source="https://github.com/NavarroLorenzo/ingsoft3-tp01"
LABEL org.opencontainers.image.description="Frontend del gestor de gastos personales"

# Configura Nginx para servir la aplicación y sus rutas.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copia únicamente los archivos estáticos generados en la etapa de build.
COPY --from=build /app/dist /usr/share/nginx/html

# Documenta el puerto HTTP expuesto por Nginx.
EXPOSE 80
