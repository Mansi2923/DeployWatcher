# DeployWatch - CI/CD Deployment Tracker Dashboard

DeployWatch is a real-time dashboard for tracking application deployments across different environments. It provides visibility into the deployment lifecycle, from queued to completed states.

## Features

- Real-time deployment status tracking
- Support for multiple environments (Dev, Staging, Prod)
- GitHub/GitLab webhook integration
- Modern, responsive UI
- Kubernetes-native deployment

## Architecture

```
+---------------------+
|    React + Redux    |
|  Deployment UI      |
+----------+----------+
           |
(fetches status via REST API)
           |
+----------v----------+
|      GoLang API     |
| Receives webhook,   |
| tracks deployments  |
+----------+----------+
           |
+----------v----------+
|  Kubernetes Cluster |
|  Apps & Deployments |
+----------+----------+
           |
+----------v----------+
|     Terraform       |
| Provisions infra,   |
| K8s resources       |
+---------------------+
```

## Prerequisites

- Go 1.21 or later
- Node.js 16 or later
- Kubernetes cluster
- Terraform 1.0 or later

## Setup

### Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

### Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

### Infrastructure

1. Navigate to the infra directory:
   ```bash
   cd infra
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

## API Endpoints

- `GET /api/deployments` - List all deployments
- `POST /api/deployments` - Create a new deployment
- `PUT /api/deployments/{id}` - Update a deployment
- `GET /api/deployments/{id}` - Get a specific deployment

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 