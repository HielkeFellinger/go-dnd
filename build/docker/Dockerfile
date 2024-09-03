FROM golang:1.22 AS builder
RUN mkdir /go-dnd
WORKDIR /go-dnd
COPY app/ app/
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main .
RUN echo $(ls -al)

FROM scratch
WORKDIR /
COPY content/ content/
COPY web/ web/
COPY campaign_data campaign_data/
COPY .env .
COPY --from=builder /go-dnd/build/main ./main
EXPOSE 8000
ENTRYPOINT ["/main"]
