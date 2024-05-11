FROM alpine
COPY invest_metrics /usr/bin/invest_metrics
ENTRYPOINT ["/usr/bin/invest_metrics"]