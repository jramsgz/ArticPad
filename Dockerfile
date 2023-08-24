# Use the latest Alpine Linux image.
FROM alpine:latest as release

# Set working directory for this stage.
WORKDIR /app

# Copy the compiled executable.
COPY  articpad .
# Copy the .env file.
COPY config/.env.sample ./config/.env
# Copy templates.
COPY templates ./templates
# Copy the static files.
COPY ui/dist ./static

# Add packages
RUN apk -U upgrade \ 
    && apk add --no-cache dumb-init curl ca-certificates tzdata

# Healthcheck
HEALTHCHECK --start-period=10s --interval=10s --timeout=5s \
  CMD curl -f http://localhost:8080/health || exit 1

# Run application and expose port 8080.
EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./articpad"]
