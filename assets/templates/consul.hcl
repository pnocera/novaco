log_level = "{{.LogLevel}}"

data_dir = "{{.DataDir}}"

datacenter = "{{.DataCenter}}"

client_addr = "{{.ClientAddr}}"

server = {{.Server}}

ui_config {
    enabled = {{.UiEnabled}}
}

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
}

service {
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


service {
  name = "webapi"
  id   = "webapi"
  port = {{.ApiPort}}
  tags = ["primary"]

  checks = [
    {
        id = "gitea-healthcheck"
        name = "Gitea Healthcheck"
        http = "http://{{.ApiHost}}:{{.ApiPort}}/api/v1/healthz"
        interval = "10s"
    }
  ]

}