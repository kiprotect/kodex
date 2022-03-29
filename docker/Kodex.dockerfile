FROM alpine:3.14.4

# Create a group and user
RUN addgroup --gid 9999 kodex && adduser --disabled-password --gecos '' --uid 9999 -G kodex -s /bin/ash kodex

WORKDIR /app
COPY bin/kodex /app
COPY entrypoint-kodex.sh /app

ENTRYPOINT ["/bin/sh", "./entrypoint-kodex.sh"]
