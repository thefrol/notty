FROM golang:1.21-alpine as base

WORKDIR /project
COPY . .

RUN ls
RUN mkdir /binaries
RUN ls cmd
RUN go build -o /binaries/ ./cmd/...

FROM scratch AS stage

ARG cmd

COPY --from=base /binaries/${cmd} ./service

ENTRYPOINT [ "./service" ]