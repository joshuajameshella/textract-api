# Contact Lambda

---

### Building the Lambda Function:

- In a terminal window, navigate to the Textract folder: `/lambda/Contact`

- Run the following command: `./build.sh`.

- (If there is a 'permission denied' error, you may need to run this command first: `chmod +x build.sh`)

- This will create a `dist/` directory, containing a `.zip` file.

- To deploy the project, open the AWS Lambda function in the AWS console, and upload the `.zip` file.

### Testing the Lambda Function:

TODO