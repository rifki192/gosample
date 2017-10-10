server {
    server_name  gosample.tokopedia.com;
    listen 80;

    access_log /var/log/nginx/gosample.access.log main;
    error_log  /var/log/nginx/gosample.error.log;

    root /var/www/gosample/public;

    location / {
        proxy_pass http://127.0.0.1:9000;
    }
}
