#! /bin/bash -x

# Working test


# Trans Image name
TRANS_IMAGE=" nwtgck/trans-server-akka:v1.13.0"
# Pull the image
docker pull $TRANS_IMAGE
# Run trans server on Docker
TRANS_CONT_ID=`docker run -d -p 8989:80 $TRANS_IMAGE`

# Build
dep ensure
go build -o trans

# Show version
./trans version
./trans --version

# Server URL setting
export TRANS_SERVER_URL=http://localhost:8989

# Create random 10MB file
dd if=/dev/urandom of=10MB.file bs=1024 count=10240

# Send 10MB.file
FILE_ID1=`./trans send 10MB.file`

# Get the file
./trans get $FILE_ID1

# Diff
diff 10MB.file $FILE_ID1

# Remove downloaded file
rm $FILE_ID1

# Parallel download
./trans get --parallel $FILE_ID1

# Remove downloaded file
rm $FILE_ID1

# Send a file with flags
FILE_ID2=`./trans send --duration=10s --get-times=3 --id-length=32 --deletable --delete-key=1234 --secure-char 10MB.file`

# Delete the file
./trans delete --delete-key=1234 $FILE_ID2


# Remove 10MB.file
rm 10MB.file

# Send by pipe
FILE_ID3=`echo "hello, world via pipe!" | ./trans send`

# Send a directory
./trans send ./vendor/

# Print the content
./trans get --stdout $FILE_ID3

# Clean up the container
docker stop $TRANS_CONT_ID
docker rm $TRANS_CONT_ID
