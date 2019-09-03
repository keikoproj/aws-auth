FROM golang:1.12-alpine AS builder

RUN apk add --update \
    python \
    python-dev \
    py-pip \
    build-base \
  && pip install awscli

WORKDIR /go/src/github.com/keikoproj/aws-auth
COPY . .
RUN make build
RUN cp ./bin/aws-auth /bin/aws-auth \
    && chmod +x /bin/aws-auth
ENV HOME /root

CMD ["/bin/aws-auth"]
