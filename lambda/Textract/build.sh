rm -rf dist
mkdir dist
env GOOS=linux go build -ldflags="-s -w" -o main .
zip Textract.zip main
mv Textract.zip ./dist/
rm main