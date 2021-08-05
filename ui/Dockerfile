FROM node AS builder
RUN mkdir /code

COPY *.json /code/
COPY webpack.config.js /code/

COPY src /code/src/

WORKDIR /code

ARG ORBITAL_SERVER_ADDRESS
ENV ORBITAL_SERVER_ADDRESS=${ORBITAL_SERVER_ADDRESS}

RUN npm install && npm run build

FROM nginx:alpine
WORKDIR /usr/share/nginx/html
# Remove default nginx static assets
RUN rm -rf ./*
# Copy static assets from builder stage
COPY --from=builder /code/dist .
# Containers run nginx with global directives and daemon off
ENTRYPOINT ["nginx", "-g", "daemon off;"]
