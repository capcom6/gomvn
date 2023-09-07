data "docker_registry_image" "app" {
  name = "capcom6/${var.app-name}:${var.app-version}"
}

data "docker_network" "proxy" {
  name = "proxy"
}


resource "docker_image" "app" {
  name          = data.docker_registry_image.app.name
  pull_triggers = [data.docker_registry_image.app.sha256_digest]
  keep_locally  = true
}

resource "docker_secret" "config" {
  name = "${var.app-name}-config.yml-${replace(timestamp(), ":", ".")}"
  data = var.app-config-b64

  lifecycle {
    ignore_changes        = [name]
    create_before_destroy = true
  }
}

resource "docker_service" "app" {
  name = var.app-name

  task_spec {
    container_spec {
      image = docker_image.app.image_id

      secrets {
        secret_id   = docker_secret.config.id
        secret_name = docker_secret.config.name
        file_name   = "/app/config.yml"
      }
    }
    networks_advanced {
      name = data.docker_network.proxy.id
    }

    resources {
      limits {
        nano_cpus    = var.cpu-limit
        memory_bytes = var.memory-limit
      }

      reservation {
        nano_cpus    = 10 * 10000000
        memory_bytes = 16 * 1024 * 1024
      }
    }
  }

  labels {
    label = "traefik.enable"
    value = true
  }
  labels {
    label = "traefik.http.routers.${var.app-name}.rule"
    value = "Host(`${var.app-host}`)"
  }
  labels {
    label = "traefik.http.routers.${var.app-name}.entrypoints"
    value = "https"
  }
  labels {
    label = "traefik.http.routers.${var.app-name}.tls"
    value = true
  }
  labels {
    label = "traefik.http.services.${var.app-name}.loadbalancer.server.port"
    value = 8080
  }
}
