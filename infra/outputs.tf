# Los outputs permiten consultar resultados con terraform output después del apply.
# Sólo se muestran datos no sensibles. La contraseña no es un output.

# El host lo devuelve Neon: no se copia de otro entorno ni se escribe a mano.
output "db_host" {
  value = neon_project.preprod.database_host
}

# Muestra el nombre de la base declarada para comprobar qué entorno se creó.
output "database" {
  value = neon_database.preprod.name
}

# Se consultan después de cada apply: Render asigna estas URLs al crear los servicios.
output "api_url" {
  value = render_web_service.api.url
}

output "front_url" {
  value = render_web_service.front.url
}
