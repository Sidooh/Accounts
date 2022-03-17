FROM golang:1.17-bullseye as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build -o /server


FROM gcr.io/distroless/base-debian11

COPY --from=build /server /server

COPY .env.example .env
# USER nonroot:nonroot
# RUN echo $USER

EXPOSE 3000

ENTRYPOINT [ "/server" ]