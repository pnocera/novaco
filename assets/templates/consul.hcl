log_level = "{{.LogLevel}}"

data_dir = "{{.DataDir}}"

datacenter = "{{.DataCenter}}"

client_addr = "{{.ClientAddr}}"

server = {{.Server}}

ui_config {
    enabled = {{.UiEnabled}}
}

auto_reload_config = true

recursors = ["1.1.1.1"]


bind_addr = "{{.BindAddr}}"

bootstrap = {{.Bootstrap}}

connect {
    enabled = {{.ConnectEnabled}}
}

addresses {
    http = "{{.AddressesHttp}}"
}

ports {
    http = {{.ConsulPort}}
    dns = 53
    grpc = 8502
}

services {
  name = "gitea"
  id   = "gitea"
  port = {{.GitPort}}
  tags = ["primary"]

  checks = [
    {
        id = "gitea-healthcheck"
        name = "Gitea Healthcheck"
        http = "http://{{.GitHost}}:{{.GitPort}}/api/healthz"
        interval = "10s"
    }
  ]

}


services {
  name = "webapi"
  id   = "webapi"
  port = {{.ApiPort}}
  tags = ["primary"]

  checks = [
    {
        id = "webapi-healthcheck"
        name = "Web API Healthcheck"
        http = "http://{{.ApiHost}}:{{.ApiPort}}/api/v1/healthz"
        interval = "10s"
    }
  ]

}