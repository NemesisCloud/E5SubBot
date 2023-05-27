export GO111MODULE=on CGO_ENABLED=0 && \
go mod download && \
go build -ldflags '-w -s' -o e5sb . && \
ls -la && pwd && \
mkdir -p /home/runner/home && \
cp e5sb /home/runner/home/e5sb
