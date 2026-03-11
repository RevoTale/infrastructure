FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY includes/ /etc/nginx/includes/
COPY templates/ /etc/nginx/templates/

RUN rm -f /etc/nginx/conf.d/default.conf && mkdir -p /var/cache/nginx
VOLUME /etc/nginx/certs/
EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
