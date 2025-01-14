ARG GOLANG_IMAGE=golang:1.22
FROM ${GOLANG_IMAGE}

WORKDIR /src
COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/server/main.go

COPY config/docker.yaml /config/docker.yaml
RUN mkdir -p /data; touch /data/app.db

EXPOSE 5051

ENTRYPOINT ["/bin/app"]