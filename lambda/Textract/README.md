# Textract Lambda

---

### Building the Lambda Function:

- In a terminal window, navigate to the Textract folder: `/lambda/Textract`

- Run the following command: `./build.sh`.

- (If there is a 'permission denied' error, you may need to run this command first: `chmod +x build.sh`)

- This will create a `dist/` directory, containing a `.zip` file.

- To deploy the project, open the AWS Lambda function in the AWS console, and upload the `.zip` file.

### Testing the Lambda Function:

Upload a document you wish to process to the Extraction Engine temporary [S3 bucket](https://s3.console.aws.amazon.com/s3/buckets/extraction-engine-temporary?region=eu-west-2&tab=objects).

Using an API testing tool (such as Postman or Insomnia), make a `POST` request to the following URL:

`https://jmkoanfttc.execute-api.eu-west-2.amazonaws.com/V1/TextractProcessing`

With the request body:

```
{
    "userID": "Josh Hellawell",
    "data": YOUR-S3-FILENAME
}
```

The process extracts the data from the S3 file location, and saves a copy of the extracted text to a `JSON` document with the same ID.
This means if a user re-submits a document, we can save money and time by returning the previously extracted `JSON` document.

After a couple of seconds, the response body should contain the following data:

```
{
    "response Data": "https://extraction-engine-temporary.s3.eu-west-2.amazonaws.com/YOUR-S3-FILENAME.json",
    "data": {
        [
            {
                "ID": "b294db35-dd6b-47b1-bf98-47d31986c61a",
                "Page": 1,
                "Word": "HELLO",
                "Confidence": 99.87960815429688
            },
        ],
        [
            {
                "ID": "b294db35-dd6b-47b1-bf98-47d31986c61a",
                "Page": 2,
                "Word": "WORLD",
                "Confidence": 99.87960815429688
            },
        ],
        ...
    }
}
```

### Next Steps

AWS Textract takes a good few seconds to process documents. Especially when fetching data from S3.
This means we're paying to run a Lambda service that is sat idle.
Using Lambda functions is cheaper at the moment, but if this project scales, it will definitely be cheaper to host on EC2.

Maybe even call Textract from the front end, and have a Lambda to just log user events.