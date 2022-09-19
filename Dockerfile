# Build golang binary in container Stage1:
FROM golang:1.16 as buld

COPY . /src

WORKDIR /src

# Build go binary
# CGO_ENABLED - tells the compiler to disable cgo and statically link any C bindings
# GOOS= linux  - creates linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs


# Scratch image contaions no distribution files. It will have only service binary Stage2:
FROM scratch

COPY --from=build /src/kvs .

COPY --from=build /src/*.pem .

EXPOSE 8080

CMD ["/kvs"]

# We need to run command `CGO_ENABLED=0 GOOD=linux go build -a -o kvs` to build our golang app binary
