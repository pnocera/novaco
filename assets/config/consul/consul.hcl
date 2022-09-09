log_level = "DEBUG"

data_dir = "e:/Projects/nomad/server_install/data/consul"

datacenter = "dc1"

client_addr = "192.168.1.145"

server = true

ui_config {
    enabled = true
}

bind_addr = "192.168.1.145"

bootstrap = true

connect {
    enabled = true
}

addresses {
    http = "192.168.1.145 127.0.0.1"
}

ports {
    http = 8500
}

service {
  name = "gitea"
  id   = "gitea"
  port = 8888
  tags = ["primary"]

  checks = [
    {
        id = "gitea-healthcheck"
        name = "Gitea Healthcheck"
        http = "http://192.168.1.145:8888/api/healthz"
        interval = "10s"
    }
  ]

}