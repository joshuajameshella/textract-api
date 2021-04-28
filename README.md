# Textract API

This is the document processing API service for Extraction Engine.

It uses AWS Lambda to process API queries, and connects to AWS Textract for document processing.

---

### Build & Deploy Project

From the project root, run the following command to build the necessary files:
`make lambda`

This will create a `build/` directory, containing a `main.zip` file.

To deploy the project, open the AWS Lambda function you wish to deploy to, and upload the `main.zip` file.

---
