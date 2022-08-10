log_level = "DEBUG"

data_dir = "e:/Projects/nomad/server_install/assets/data/consul"

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