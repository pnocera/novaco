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

service {
  name = "gitea"
  id   = "gitea"
  port = 8888
  tags = ["primary"]

  checks = [
    {
        id = "gitea-healthcheck"
        name = "Gitea Healthcheck"
        http = "http://{{.BindAddr}}:8888/api/healthz"
        interval = "10s"
    }
  ]

}