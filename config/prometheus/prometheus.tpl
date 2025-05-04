global:
  scrape_interval:     5s
  evaluation_interval: 5s

scrape_configs:
  - job_name: 'prometheus'

    scrape_interval: 800ms

    static_configs:
      - targets: ['{{.Metrics.Host)}:{{.Metrics.Port}}']

  - job_name: 'node-exporter'

    scrape_interval: 5s

    static_configs:
      - targets: ['{{.NodeExporter.Host}}:{{.NodeExporter.Port}}']