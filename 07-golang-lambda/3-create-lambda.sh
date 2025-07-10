# Default values
FUNCTION_NAME="go-lambda"
ZIP_FILE="myFunction.zip"
HANDLER="bootstrap"
RUNTIME="provided.al2023"
ROLE_ARN="arn:aws:iam::028155753870:role/lambda-ex"

# Create the Lambda function
aws lambda create-function \
    --function-name "$FUNCTION_NAME" \
    --zip-file "fileb://$ZIP_FILE" \
    --architectures arm64 \
    --handler "$HANDLER" \
    --runtime "$RUNTIME" \
    --role "$ROLE_ARN"