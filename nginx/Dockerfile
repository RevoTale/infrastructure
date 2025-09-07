FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /var/cache/nginx
VOLUME /etc/nginx/certs/
EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
