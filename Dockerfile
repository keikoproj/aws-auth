FROM golang:1.24-alpine as build

RUN apk add --update --no-cache \
    curl \
    build-base \
    git \
    python3 \
    py3-pip
# Install AWS CLI using Alpine package manager or setup a virtual environment
RUN python3 -m venv /tmp/venv && \
    . /tmp/venv/bin/activate && \
    pip install awscli && \
    cp -r /tmp/venv/bin/aws* /usr/local/bin/

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
