log_level = "{{.LogLevel}}"

data_dir = "{{.DataDir}}"

server {
    enabled = true
    bootstrap_expect = {{.BootstrapExpect}}
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