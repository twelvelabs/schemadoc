FROM alpine:3.21

COPY schemadoc /usr/local/bin/schemadoc

CMD ["/usr/local/bin/schemadoc"]
