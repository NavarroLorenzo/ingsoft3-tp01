# Define la versión mínima de Terraform y las versiones exactas de sus conectores.
# Sin un bloque backend, el estado se guarda localmente en infra/terraform.tfstate.
terraform {
  required_version = ">= 1.5.0"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0"
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
  }
}

# Cada provider toma las credenciales de la terminal; no se escriben en este archivo.
# Neon lee NEON_API_KEY. Render lee RENDER_API_KEY y RENDER_OWNER_ID.
provider "neon" {}
provider "render" {}

# Un proyecto propio permite destruir preprod sin afectar las bases existentes.
resource "neon_project" "preprod" {
  name = "ingsoft3-preprod"
  # La API de Neon necesita este dato cuando se usa una clave personal.
  org_id = var.neon_org_id
  # Se solicitan seis horas de historial, como indica la guía para el plan gratuito.
  history_retention_seconds = 21600
}

# Crea el usuario de la base en la rama principal del proyecto recién creado.
# Neon genera su contraseña; Terraform la conserva como dato sensible en el estado.
resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "app_preprod_owner"
}

# Crea la base y la asigna al usuario anterior. Las referencias establecen el orden:
# primero el proyecto, después el usuario y finalmente la base.
resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "app_preprod"
  owner_name = neon_role.preprod.name
}

# Crea una API nueva desde la imagen publicada del TP7; no compila ni toca QA/PROD.
resource "render_web_service" "api" {
  name   = "ingsoft3-api-preprod"
  plan   = "free"
  region = var.region

  # La API no tiene una ruta en /; /health comprueba también la conexión a la base.
  health_check_path = "/health"

  runtime_source = {
    image = {
      image_url = var.image_repo_api
      tag       = var.image_tag
    }
  }

  # Nuestra app en Go arma la conexión a partir de DB_*, en lugar de recibir
  # ConnectionStrings__Default como la app .NET de ejemplo de la guía.
  # Los datos propios de la base salen de los recursos anteriores, sin copiarlos.
  env_vars = {
    DB_HOST     = { value = neon_project.preprod.database_host }
    DB_PORT     = { value = "5432" }
    DB_USER     = { value = neon_role.preprod.name }
    DB_PASSWORD = { value = neon_role.preprod.password }
    DB_NAME     = { value = neon_database.preprod.name }
    DB_SSLMODE  = { value = "require" }
    SERVER_PORT = { value = "8080" }

    # La app necesita una clave para firmar los tokens de login. Render la genera
    # al crear el servicio; su valor no se escribe en los archivos versionados.
    JWT_SECRET = { generate_value = true }
  }
}

# Crea el frontend de preprod desde la misma versión del TP7 que usa la API.
resource "render_web_service" "front" {
  name   = "ingsoft3-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  # Nginx recibe la URL de SU API al arrancar. Esta referencia también hace que
  # Terraform cree primero la API y destruya primero el frontend al bajar preprod.
  env_vars = {
    BACKEND_URL  = { value = render_web_service.api.url }
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}
