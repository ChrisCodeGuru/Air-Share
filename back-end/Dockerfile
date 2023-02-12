# pull official base image
FROM golang:alpine

# set working directory
WORKDIR /app

# add app
COPY . ./

# start app
CMD ["go", "run", "main.go"]