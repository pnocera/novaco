storage "consul" {
  address = "{{.StorageConsulAddress}}"
  path    = "{{.StorageConsulPath}}"
}

log_level = "{{.LogLevel}}"

listener "tcp" {
  address     = "{{.TcpAddress}}"
  tls_disable = {{.TcpTlsDisable}}
}

service_registration "consul" {
  address      = "{{.StorageConsulAddress}}"
}

// telemetry {
//   statsite_address = "{{.TelemetryStatsdAddress}}"
//   disable_hostname = {{.TelemetryDisableHostname}}
// }

ui = {{.UiEnabled}}