global:
  scrape_interval:     5s
  evaluation_interval: 5s

scrape_configs:
  - job_name: 'prometheus'

    scrape_interval: 800ms

    static_configs:
      - targets: ['{{.metrics.host}}:{{.metrics.port}}']