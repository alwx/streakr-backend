FROM golang:latest as build
COPY . /go/src/streakr-backend/
WORKDIR /go/src/streakr-backend/
RUN go get ./cmd/streakr-backend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/streakr-server ./cmd/streakr-go

FROM scratch
COPY --from=build /go/src/streakr-backend/out/streakr-server /bin/streakr-go
COPY --from=build /go/src/streakr-backend/configs /etc/streakr-go/configs
COPY --from=build /go/src/streakr-backend/migrations /etc/streakr-go/migrations
ENTRYPOINT ["./bin/streakr-go"]