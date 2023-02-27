# Use the latest Alpine Linux image.
FROM alpine:latest as release

# Set working directory for this stage.
WORKDIR /app

# Copy our compiled executable.
COPY  articpad .
# Copy our .env file.
COPY config/.env.sample ./config/.env
# Copy our static files.
COPY ui/dist ./static

# Add packages
RUN apk -U upgrade \ 
    && apk add --no-cache dumb-init curl ca-certificates tzdata

# Healthcheck
HEALTHCHECK --start-period=10s --interval=10s --timeout=5s \
  CMD curl -f http://localhost:3000/health || exit 1

# Run application and expose port 3000.
EXPOSE 3000
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./articpad"]
