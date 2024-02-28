FROM alpine:3.19.1 AS release

# Add nonroot user and group.
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Set working directory for this stage.
WORKDIR /app

# Copy the compiled executable.
COPY articpad .
# Copy the .env file.
COPY config/.env.sample ./config/.env
# Copy templates.
COPY templates/ ./templates/
# Copy locales.
COPY locales/ ./locales/
# Copy the static files.
COPY ui/dist/ ./static/

# Add packages
RUN apk -U upgrade \ 
    && apk add --no-cache dumb-init curl ca-certificates tzdata

# Healthcheck
HEALTHCHECK --start-period=10s --interval=10s --timeout=5s \
  CMD curl -f http://localhost:8080/health || exit 1

# Set the nonroot user as the default user.
USER nonroot

# Run application and expose port 8080.
EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./articpad"]
