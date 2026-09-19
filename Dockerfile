FROM amd64/alpine
EXPOSE 8080
WORKDIR /app
COPY short2 .
CMD ["/app/short2"]
