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