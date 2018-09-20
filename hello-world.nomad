job "hello-world" {
  datacenters = ["dc1"]
  type = "service"
  
  group "grp1" {
    count = 1

    task "run-hello-world" {
      driver = "docker"
      config {
        image = "toomyem/hello-world:0.2"
        port_map = {
          http = 9000
        }
      }

      service {
        name = "hello-world"
        port = "http"
        check {
          type = "http"
          path = "/health"
          interval = "10s"
          timeout = "2s"
        }
      }

      resources {
        memory = 100 # MB

        network {
          port "http" {}
        }
      }
    }
  }
}

