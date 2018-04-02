#
# Expect to setup environment variables TITLE and VERSION
# Pass integer value 1-65535 for listening port
# Default port is 8282
#

FROM scratch

#   nobody:nobody
USER 65534:65534

COPY server /

ENTRYPOINT ["/server"]