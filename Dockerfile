FROM reg.c5h.io/golang as builder

WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -trimpath -ldflags=-buildid= -o main ./

FROM reg.c5h.io/base

COPY --from=builder /app/main /todoistager
CMD ["/todoistager"]