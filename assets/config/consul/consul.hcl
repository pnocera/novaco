log_level = "DEBUG"

data_dir = "e:/Projects/nomad/server_install/data/consul"

datacenter = "dc1"

client_addr = "192.168.1.145"

server = true

ui_config {
    enabled = true
}

auto_reload_config = true

recursors = ["1.1.1.1"]


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
    dns = 53
    grpc = 8502
}

services {
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


services {
  name = "webapi"
  id   = "webapi"
  port = 7788
  tags = ["primary"]

  checks = [
    {
        id = "webapi-healthcheck"
        name = "Web API Healthcheck"
        http = "http://192.168.1.145:7788/api/v1/healthz"
        interval = "10s"
    }
  ]

}