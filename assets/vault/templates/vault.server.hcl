storage "consul" {
  address = "{{.StorageConsulAddress}}"
  path    = "{{.StorageConsulPath}}"
}

listener "tcp" {
  address     = "{{.TcpAddress}}"
  tls_disable = {{.TcpTlsDisable}}
}

telemetry {
  statsite_address = "{{.TelemetryStatsdAddress}}"
  disable_hostname = {{.TelemetryDisableHostname}}
}

ui = {{.UiEnabled}}