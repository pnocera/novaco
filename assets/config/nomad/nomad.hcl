log_level = "DEBUG"

data_dir = "e:/Projects/nomad/server_install/data/nomad"

server {
    enabled = true
    bootstrap_expect = 1
}

client {
    enabled = true
    servers = ["0.0.0.0"]
}

plugin "raw_exec" {
    config {
        enabled = true
    }
}

advertise {
    http = "192.168.1.145"
    rpc  = "192.168.1.145"
    serf = "192.168.1.145"
}

addresses {
    http = "0.0.0.0"
    rpc  = "192.168.1.145"
    serf = "192.168.1.145"
}