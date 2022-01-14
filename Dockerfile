MAINTAINER suwei007@gmail.com
FROM scratch
# Copy our static executable.
COPY ./server/server /server
# Run the hello binary.
EXPOSE 50051
ENTRYPOINT ["/server"]