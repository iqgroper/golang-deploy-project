FROM golang:1.16-alpine AS build_stage
COPY ./crudapp /go/src/crudapp
WORKDIR /go/src/crudapp
RUN go install .


FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/bin/crudapp /app_binary/
RUN chmod +x ./crudapp
EXPOSE 8080/tcp
ENTRYPOINT ./crudapp

# docker build --tag=crudapp:latest .
# docker run --rm -p 8080:8080 crudapp