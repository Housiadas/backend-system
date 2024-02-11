ARG VERSION
ARG UID
ARG GID
ARG LINKER_FLAGS

# Build the Go Binary.
FROM golang:$VERSION as banking-api

# Copy the source code into the container.
COPY . /api

# Build the app binary.
WORKDIR /api/app/
RUN go build -ldflags=$LINKER_FLAGS -o=./banking-api

# Run the Go Binary in Alpine.
FROM alpine:3.19

COPY --from=banking-api --chown=$UID:$GID ./api/app/banking-api ./api/banking-api
COPY --from=banking-api --chown=$UID:$GID ./api/app.env.example ./api/app.env

WORKDIR /api

USER $UID:$GID

CMD ["./banking-api"]
