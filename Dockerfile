FROM scratch
MAINTAINER suwei007@gmail.com
# Copy our static executable.
COPY ./code/server/server /server
# Run the hello binary.
EXPOSE 50051
ENTRYPOINT ["/server"]