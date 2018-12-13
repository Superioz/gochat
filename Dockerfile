#base image
FROM alpine

#port
EXPOSE 6000

#copy executable from host to container
COPY server.exe /etc/chatserver

#execute binary from container
CMD ./etc/chatserver
