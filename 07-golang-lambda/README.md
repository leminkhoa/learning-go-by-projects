# AWS Lambda Go Example

A minimal AWS Lambda function written in Go, demonstrating how to build, package, and deploy a serverless function to AWS Lambda using the custom runtime. This project is ideal for learning or as a template for your own Go-based Lambda functions.

---

## 🚀 Features
- Written in idiomatic Go
- AWS Lambda custom runtime (arm64, provided.al2023)
- Simple event/response structure
- Automated scripts for IAM role creation, build, packaging, and deployment

---

## 📦 Project Structure
```
├── 1-create-role.sh         # Script to create IAM role for Lambda
├── 2-create-zip.sh          # Script to build and package the Go binary
├── 3-create-lambda.sh       # Script to deploy the Lambda function
├── bootstrap                # Compiled Go binary (output)
├── go.mod / go.sum          # Go module files
├── main.go                  # Lambda function source code
├── myFunction.zip           # Deployment package (output)
├── trust-policy.json        # IAM trust policy for Lambda
```

---

## 🛠️ Prerequisites
- [Go 1.24+](https://golang.org/dl/)
- [AWS CLI](https://aws.amazon.com/cli/) configured with appropriate permissions
- AWS account with permissions to create IAM roles and Lambda functions

---

## ⚙️ Setup & Deployment

### 1. Create IAM Role
This role allows Lambda to assume execution permissions.
```bash
bash 1-create-role.sh
```
- Uses `trust-policy.json` to define trust relationship
- Attaches `AWSLambdaBasicExecutionRole` policy

### 2. Build & Package the Function
Compile the Go code for AWS Lambda (Linux/arm64) and package it as a zip file.
```bash
bash 2-create-zip.sh
```
- Produces `bootstrap` binary and `myFunction.zip`

### 3. Deploy to AWS Lambda
Create the Lambda function using the packaged zip and IAM role.
```bash
bash 3-create-lambda.sh
```
- Uses `provided.al2023` runtime and arm64 architecture

---

## 📝 Example: Event & Response

**Event JSON:**
```json
{
  "what is your name?": "Alice",
  "How old are you": 30
}
```

**Response JSON:**
```json
{
  "Answer:": "Alice is 30 years old"
}
```

---

## 📚 References & Acknowledgments
- [AWS Lambda Go Documentation](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)
- [aws/aws-lambda-go](https://github.com/aws/aws-lambda-go)

---

## 📝 License
This project is for educational purposes. Add your license here if needed. 