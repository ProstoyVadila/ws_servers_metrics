FROM rust:1.77.1-alpine3.19 as builder
# This is important, see https://github.com/rust-lang/docker-rust/issues/85
ENV RUSTFLAGS="-C target-feature=-crt-static"
# ENV ROCKET_ENV=production
# ENV ROCKET_PORT=8082

# if needed, add additional dependencies here
RUN apk add --no-cache musl-dev pkgconfig openssl-dev
# set the workdir and copy the source into it
WORKDIR /app
COPY ./ /app
# do a release build
RUN cargo build --release
RUN strip target/release/rust_ws_app

# use a plain alpine image, the alpine version needs to match the builder
FROM alpine:3.19

# ENV RUST_LOG=info
# ENV ROCKET_PORT=8080

# if needed, install additional dependencies here
RUN apk add --no-cache libgcc curl
# copy the binary into the final image
# COPY --from=builder /app/Rocket.toml .
COPY --from=builder /app/target/release/rust_ws_app .

ENV ROCKET_PORT=8000
ENV ROCKET_ADDRESS=0.0.0.0
ENV ROCKET_WORKERS=4

# set the binary as entrypoint
ENTRYPOINT ["./rust_ws_app"]
