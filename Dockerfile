# syntax=docker/dockerfile:1

# Frontend builder stage
FROM --platform=$BUILDPLATFORM node:22.13.1-alpine AS frontendbuilder

WORKDIR /build

ENV PNPM_CACHE_FOLDER=.cache/pnpm/ \
    PUPPETEER_SKIP_DOWNLOAD=true \
    CYPRESS_INSTALL_BINARY=0

# Copy package files first to leverage layer caching
COPY frontend/package.json frontend/pnpm-lock.yaml ./
COPY frontend/patches ./patches/

# Install dependencies in a single layer
RUN npm install -g corepack && \
    corepack enable && \
    pnpm install --frozen-lockfile

# Copy source files after dependencies
COPY frontend/ ./

# Build the frontend
RUN pnpm run build

# API builder stage
FROM --platform=$BUILDPLATFORM ghcr.io/techknowlogick/xgo:go-1.23.x AS apibuilder

# Install mage in a separate layer
RUN go install github.com/magefile/mage@latest && \
    mv /go/bin/mage /usr/local/go/bin

WORKDIR /go/src/code.vikunja.io/api

# Copy go.mod and go.sum first to leverage layer caching
COPY go.mod go.sum ./

# Download dependencies in a separate layer
ENV GOPROXY=https://goproxy.kolaente.de
RUN go mod download

# Copy the rest of the source code
COPY . ./
COPY --from=frontendbuilder /build/dist ./frontend/dist

ARG TARGETOS TARGETARCH TARGETVARIANT RELEASE_VERSION
ENV RELEASE_VERSION=$RELEASE_VERSION

# Build the application
RUN export PATH=$PATH:$GOPATH/bin && \
    mage build:clean && \
    mage release:xgo "${TARGETOS}/${TARGETARCH}/${TARGETVARIANT}"

# Final stage
FROM scratch

# Add metadata labels in a single layer
LABEL org.opencontainers.image.authors='maintainers@vikunja.io' \
      org.opencontainers.image.url='https://vikunja.io' \
      org.opencontainers.image.documentation='https://vikunja.io/docs' \
      org.opencontainers.image.source='https://code.vikunja.io/vikunja' \
      org.opencontainers.image.licenses='AGPLv3' \
      org.opencontainers.image.title='Vikunja'

WORKDIR /app/vikunja
ENV VIKUNJA_SERVICE_ROOTPATH=/app/vikunja/ \
    VIKUNJA_DATABASE_PATH=/db/vikunja.db

COPY --from=apibuilder /build/vikunja-* vikunja
COPY --from=apibuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1000
EXPOSE 3456
ENTRYPOINT [ "/app/vikunja/vikunja" ]
