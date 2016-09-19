FROM golang:1.7-alpine

ENV PROJECT	$GOPATH/src/github.com/zabawaba99/stash-commit-status-resource/
COPY . $PROJECT

WORKDIR $PROJECT
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o check ./cmd/check
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o in ./cmd/in
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o out ./cmd/out

CMD ["tar", "-cf", "-", "check", "in", "out"]
