FROM scratch

#   nobody:nobody
USER 65534:65534

ENV TITLE
ENV VERSION

EXPOSE 8282

COPY simple-test-app /

ENTRYPOINT ["/simple-test-app"]