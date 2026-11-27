# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS backend-build
WORKDIR /app
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /calculator .

FROM scratch AS backend
COPY --from=backend-build /calculator /calculator
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/calculator"]

FROM node:24-bookworm-slim AS frontend-build
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM nginx:1.30-alpine AS frontend
COPY --from=frontend-build /app/dist/ /usr/share/nginx/html/
COPY <<EOF /etc/nginx/conf.d/default.conf
server {
    listen 80;
    root /usr/share/nginx/html;

    location /api/ {
        # Go enforces the request body limit and returns the JSON error.
        client_max_body_size 0;
        proxy_request_buffering off;
        proxy_http_version 1.1;
        proxy_pass http://backend:8080;
    }
}
EOF
