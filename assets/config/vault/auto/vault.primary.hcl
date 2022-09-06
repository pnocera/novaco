storage "consul" {
  address = "192.168.1.145:8500"
  path    = "vault/"
}

log_level = "DEBUG"

listener "tcp" {
  address     = "192.168.1.145:8200"
  tls_disable = 1
}

service_registration "consul" {
  address      = "192.168.1.145:8500"
}

// telemetry {
//   statsite_address = "192.168.1.145:8125"
//   disable_hostname = true
// }

ui = true