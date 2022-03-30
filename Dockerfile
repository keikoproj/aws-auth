FROM golang:1.17-alpine as build

RUN apk add --update --no-cache \
    curl \
    py-pip \
    build-base \
    git \
  && pip install awscli

WORKDIR /go/src/github.com/keikoproj/aws-auth
COPY . .
RUN git rev-parse HEAD
RUN date +%FT%T%z
RUN make build
RUN chmod +x ./bin/aws-auth

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
COPY --from=build /go/src/github.com/keikoproj/aws-auth/bin/aws-auth /bin/aws-auth

ENV HOME /root
ENTRYPOINT ["/bin/aws-auth"]
CMD ["help"]
