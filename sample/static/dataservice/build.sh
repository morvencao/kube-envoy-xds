CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build --no-cache -f Dockerfile-dataservice -t morvencao/dataservice:v2.0 .
rm -rf main
