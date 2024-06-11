FROM golang:1.22.3 AS build
WORKDIR /src
# Copy dependencies list
COPY go.mod go.sum ./
# Build with optional lambda.norpc tag
COPY *.go .
RUN go build -tags lambda.norpc -o main

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
RUN dnf -y upgrade && dnf clean all
COPY --from=build /src/main ./main
HEALTHCHECK NONE
USER 65534
ENTRYPOINT [ "./main" ]