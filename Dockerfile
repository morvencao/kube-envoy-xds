FROM alpine:latest

ADD kube-envoy-xds /kube-envoy-xds
ENTRYPOINT ["./kube-envoy-xds"]
