# Envoy xDS backed by Kubernetes

Simple golang-based implementation of an API server that
aimed to run in Kubernetes cluster that implements the Envoy
discovery service APIs defined in [data-plane-api](https://github.com/envoyproxy/data-plane-api).

## Requirements

1. Go 1.9+
2. Docker
3. Kubernetes

## Build

1. Setup tools and dependencies

```sh
make tools
make depend.install
```

2. Generate proto files (if you update the [data-plane-api](https://github.com/envoyproxy/data-plane-api)
dependency)

```sh
make generate
```

3. Format, vet and lint the code

```sh
make check
```

5. Build

```sh
make build
```

6. Build docker image
```sh
make docker
```

7. Push docker image
```sh
make push
```

8. Clean the build
```sh
make clean
```
