# Identifica la organización de Neon donde se creará el proyecto de preprod.
# Se carga desde TF_VAR_neon_org_id en la terminal; no es una clave de API.
# Al no tener default, hay que proporcionar el ID al ejecutar plan o apply.
variable "neon_org_id" {
  description = "ID de la organización de Neon que será propietaria de preprod."
  type        = string
}

# Las imágenes son públicas y ya fueron construidas por el pipeline del TP7.
# El repositorio va sin etiqueta: Render recibe la etiqueta en un campo separado.
variable "image_repo_api" {
  description = "Repositorio público de la imagen del backend."
  type        = string
  default     = "ghcr.io/navarrolorenzo/ingsoft3-tp01-backend"
}

variable "image_repo_front" {
  description = "Repositorio público de la imagen del frontend."
  type        = string
  default     = "ghcr.io/navarrolorenzo/ingsoft3-tp01-frontend"
}

# Fija el merge del TP7 que coincide con v7.0.0; no depende de latest.
# Cambiarla después requiere recrear los servicios por la limitación del provider.
variable "image_tag" {
  description = "Etiqueta de las imágenes publicadas al cerrar el TP7."
  type        = string
  default     = "sha-fdb86c805bfcd2adcfb19febdbab22e92fba5859"
}

# Ambos servicios de Render usan Oregon, una de las regiones admitidas por el provider.
variable "region" {
  description = "Región de Render para los dos servicios de preprod."
  type        = string
  default     = "oregon"
}
