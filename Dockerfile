# Use 'scratch' image for super-mini build.
FROM scratch AS prod

# Set working directory for this stage.
WORKDIR /production

# Copy our compiled executable.
COPY  articpad .
# Copy our .env file.
COPY .env.sample ./config/.env
# Copy our static files.
COPY ./ui/dist ./static

# Run application and expose port 3000.
EXPOSE 3000
CMD ["./articpad"]
