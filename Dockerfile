FROM golang:1.20 as BUILDER

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -o main main.go 

FROM golang:1.20

WORKDIR /usr/src/app

COPY --from=BUILDER /usr/src/app/main .

EXPOSE 9001

HEALTHCHECK --timeout=2s --start-period=5s --retries=3 --interval=5s \
    CMD curl --fail http://localhost:9001/live || exit 1

ENTRYPOINT [ "/usr/src/app/main" ]
