# Extraction Engine API

This is the document processing API service for Extraction Engine.

It uses AWS Lambda to process API queries, and connects to AWS Textract for document processing.

The different Lambda functions can be found in the `lambda/` directory.

---

### Build & Deploy Project

for each lambda function, navigate into the directory, and run the following command: `build.sh`

(If there is a 'permission denied' error, you may need to run this command first: `chmod +x build.sh`)

This will create a `dist/` directory, containing a `.zip` file.

To deploy the project, open the AWS Lambda function in the AWS console, and upload the `.zip` file.

---
