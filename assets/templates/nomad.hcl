log_level = "{{.LogLevel}}"

data_dir = "{{.DataDir}}"

server {
    enabled = true
    bootstrap_expect = {{.BootstrapExpect}}
}

client {
    enabled = true
    servers = ["{{.BindAddr}}"]
}

plugin "raw_exec" {
    config {
        enabled = true
    }
}

advertise {
    http = "{{.AdvertiseAddr}}"
    rpc  = "{{.AdvertiseAddr}}"
    serf = "{{.AdvertiseAddr}}"
}

addresses {
    http = "{{.BindAddr}}"
    rpc  = "{{.AdvertiseAddr}}"
    serf = "{{.AdvertiseAddr}}"
}