#! /bin/bash -x

if [ "$TRAVIS" = "true" ]; then
  # == Move build directory ==
  cd /
  mkdir $GOPATH/src
  mv $TRAVIS_BUILD_DIR $GOPATH/src/trans-cli-go
  cd $GOPATH/src/trans-cli-go
fi

# Working test


# Trans Image name
TRANS_IMAGE=" nwtgck/trans-server-akka:v1.9.0"
# Pull the image
docker pull $TRANS_IMAGE
# Run trans server on Docker
TRANS_CONT_ID=`docker run -d -p 8080:80 $TRANS_IMAGE`

# Build
dep ensure
go build -o trans

# Server URL setting
export TRANS_SERVER_URL=http://localhost:8080

# Create random 10MB file
dd if=/dev/urandom of=10MB.file bs=1024 count=10240

# Send 10MB.file
FILE_ID=`./trans send 10MB.file`

# Get the file
./trans get $FILE_ID

# Diff
diff 10MB.file $FILE_ID

# Remove downloaded file
rm $FILE_ID

# Remove 10MB.file
rm 10MB.file

# Clean up the container
docker stop $TRANS_CONT_ID
docker rm $TRANS_CONT_ID
