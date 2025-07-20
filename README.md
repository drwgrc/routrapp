# RoutrApp

A modern web application built with Go/Gin backend and Next.js frontend.

## Project Overview

RoutrApp is a full-stack web application featuring:

- **Backend**: Gin-based REST API with Go
- **Frontend**: Next.js with React, TypeScript, and Tailwind CSS
- **Modern Stack**: Built with the latest technologies for optimal performance and developer experience

## Project Structure

```
routrapp/
├── backend/
│   └── cmd/api/main.go    # Gin application entry point
├── frontend/
│   ├── src/
│   │   └── app/         # Next.js app directory
│   ├── package.json     # Frontend dependencies
│   └── next.config.ts   # Next.js configuration
└── README.md
```

## Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** (v18 or higher)
- **Go** (v1.16 or higher)

## Setup Instructions

### Backend Setup

1. **Navigate to the backend directory:**

```bash
cd backend
```

2. **Install Go dependencies:**

   ```bash
   go mod tidy
   ```

3. **Run the backend server:**

   ```bash
   go run cmd/api/main.go
   ```

   The API will be available at `http://localhost:8080`

   - API documentation: `http://localhost:8080/swagger/index.html`

### Frontend Setup

1. **Navigate to the frontend directory:**

```bash
cd frontend
```

2. **Install dependencies:**

   ```bash
   npm install
   ```

3. **Run the development server:**

   ```bash
   npm run dev
   ```

   The frontend will be available at `http://localhost:3000`

## Available Scripts

### Backend

- `go run cmd/api/main.go` - Start development server
- `go test ./...` - Run all tests

### Frontend

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## API Endpoints

The backend currently provides these endpoints:

- `GET /` - Returns a simple "Hello World" message
- `GET /items/{item_id}` - Returns item information with optional query parameter

## Development

### Backend Development

- The backend uses Gin for rapid API development
- Swagger documentation is available at `/swagger/index.html`
- Configuration through environment variables

### Frontend Development

- Built with Next.js and React
- Uses TypeScript for type safety
- Styled with Tailwind CSS
- ESLint configured for code quality

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test both backend and frontend
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).
