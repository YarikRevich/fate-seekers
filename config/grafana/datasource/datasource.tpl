apiVersion: 1

deleteDatasources:
  - name: Default
    orgId: 1

datasources:
  - name: Default
    type: prometheus
    access: proxy
    orgId: 1
    url: http://{{.Prometheus.Host}}:{{.Prometheus.Port}}
    isDefault: true
    jsonData:
      graphiteVersion: "1.1"
      tlsAuth: false
      tlsAuthWithCACert: false
    version: 1
    editable: false


