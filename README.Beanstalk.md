# Deploying to AWS Elastic Beanstalk

This document provides instructions for deploying the Go Store application to AWS Elastic Beanstalk using Amazon Linux.

## Prerequisites

- AWS Account
- AWS CLI installed and configured
- Elastic Beanstalk CLI (eb-cli) installed

## Files Added for Elastic Beanstalk Deployment

1. **Procfile**: Specifies the command to start the application
   ```
   web: ./bin/server
   ```

2. **Port Configuration**: The application has been modified to listen on the port specified by the `PORT` environment variable, which is set by Elastic Beanstalk. If not set, it defaults to port 8080.

## Deployment Steps

1. **Build the application**:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd
   ```

2. **Create a deployment package**:
   Create a zip file containing:
   - bin/server (compiled binary)
   - templates/ (HTML templates)
   - .env (environment variables)
   - Procfile

   ```bash
   zip -r deploy.zip bin/server templates/ .env Procfile
   ```

3. **Initialize Elastic Beanstalk** (if not already done):
   ```bash
   eb init -p go go-store
   ```

4. **Create an environment**:
   ```bash
   eb create go-store-env
   ```

5. **Deploy the application**:
   ```bash
   eb deploy
   ```

## Environment Variables

Ensure the following environment variables are configured in the Elastic Beanstalk environment:

- `DB_HOST`: Database host
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password

You can set these in the Elastic Beanstalk console under Configuration > Software > Environment properties.

## Troubleshooting

- **Application not starting**: Check the Elastic Beanstalk logs for errors
- **Database connection issues**: Verify that the database security group allows connections from the Elastic Beanstalk environment
- **Port issues**: Ensure the application is listening on the port specified by the `PORT` environment variable