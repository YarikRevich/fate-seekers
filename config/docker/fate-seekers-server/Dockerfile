FROM golang:1.22.10-alpine as build

LABEL AUTHOR="YarikRevich"
LABEL VERSION="1.0.0"

WORKDIR /app/fate-seekers-server

COPY . .

RUN make build-server && 

FROM golang:1.22.10-alpine

COPY --from=build /app/fate-seekers-server/build/fate-seekers-server ./build/fate-seekers-server

