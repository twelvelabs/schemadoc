FROM alpine:3.23

COPY schemadoc /usr/local/bin/schemadoc

CMD ["/usr/local/bin/schemadoc"]
