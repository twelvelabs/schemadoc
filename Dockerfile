FROM alpine:3.22

COPY schemadoc /usr/local/bin/schemadoc

CMD ["/usr/local/bin/schemadoc"]
