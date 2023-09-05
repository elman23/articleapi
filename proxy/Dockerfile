FROM nginx:1.19.8-alpine

# Set TZ (TimeZone) to
# RUN apk add --no-cache tzdata # Note: if you are using an alpine based image you have to install the tzdata first
ENV TZ=Europe/Rome
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Copy server configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose ports
EXPOSE 80