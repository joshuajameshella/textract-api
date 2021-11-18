rm -rf dist
mkdir dist
env GOOS=linux go build -ldflags="-s -w" -o main .
zip Contact.zip main
mv Contact.zip ./dist/
rm main