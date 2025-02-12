spec:
  # task 1: generate HTTP requests for application URL
  # collect Iter8's built-in HTTP latency and error-related metrics
  - task: http
    with:
      duration: 2s
      errorRanges:
      - lower: 500
      url: {{ .URL }}
  # task 2: validate service level objectives for app using
  # the metrics collected in the above task
  - task: assess
    with:
      SLOs:
        upper:
        - metric: "http/error-rate"
          limit: 0
        - metric: "http/latency-mean"
          limit: 500
        - metric: "http/latency-p50"
          limit: 1000
        - metric: "http/latency-p50.0"
          limit: 1
        - metric: "http/latency-p95.0"
          limit: 2500
        - metric: "http/latency-p99"
          limit: 5000
  # tasks 3 & 4: print if SLOs are satisfied or not
  - if: SLOs()
    run: echo "SLOs satisfied"
  - if: not SLOs()
    run: echo "SLOs not satisfied"
