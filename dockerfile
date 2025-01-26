# Dockerfile for PostgreSQL
FROM postgres:latest

# Optional: Set environment variables for initial database setup
ENV POSTGRES_DB=minfin-medialist
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=root

# Optional: Add custom initialization scripts
# COPY init.sql /docker-entrypoint-initdb.d/

# Expose PostgreSQL default port
EXPOSE 5432
