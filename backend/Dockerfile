# ---- Etapa 1: build (imagen con SDK) ----
FROM golang:1.26-alpine AS build
WORKDIR /src

# Se copian primero los archivos de dependencias
# antes de descargarlas — capa cacheable: solo se rehace si cambian go.mod/go.sum.
COPY go.mod go.sum ./
RUN go mod download

# Recién ahora se copia el resto del código
COPY . .

# Compila la aplicación en modo producción
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/publish/gestor-gastos-api ./cmd/api


# ---- Etapa 2: runtime (imagen mínima, sin SDK) ----
FROM alpine:3.22 AS final
WORKDIR /app

COPY --from=build /app/publish .

EXPOSE 8080

ENTRYPOINT ["./gestor-gastos-api"]
