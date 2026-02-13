FROM debian:stable-slim


COPY web /bin/web

CMD ["/bin/web" ]