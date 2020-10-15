FROM golang:1.13-alpine

RUN apk add --update --no-cache \
    python \
    python-dev \
    py-pip \
    build-base \
    git \
  && pip install awscli

WORKDIR /go/src/github.com/keikoproj/aws-auth
COPY . .
RUN git rev-parse HEAD
RUN date +%FT%T%z
RUN make build
RUN cp ./bin/aws-auth /bin/aws-auth \
    && chmod +x /bin/aws-auth
ENV HOME /root

ENTRYPOINT ["/bin/aws-auth"]
CMD ["help"]
