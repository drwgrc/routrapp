# RoutRapp

A modern web application built with FastAPI backend and Next.js frontend.

## Project Overview

RoutRapp is a full-stack web application featuring:

- **Backend**: FastAPI-based REST API with Python
- **Frontend**: Next.js15ith React 19peScript, and Tailwind CSS
- **Modern Stack**: Built with the latest technologies for optimal performance and developer experience

## Project Structure

```
routrapp/
├── backend/
│   └── main.py          # FastAPI application entry point
├── frontend/
│   ├── src/
│   │   └── app/         # Next.js app directory
│   ├── package.json     # Frontend dependencies
│   └── next.config.ts   # Next.js configuration
└── README.md
```

## Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** (v18higher)
- **Python** (v30.8or higher)
- **pip** (Python package manager)

## Setup Instructions

### Backend Setup1 **Navigate to the backend directory:**

```bash
cd backend
```

2. **Create a virtual environment (recommended):**

   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

   3 **Install FastAPI and dependencies:**

   ```bash
   pip install fastapi uvicorn
   ```

3. **Run the backend server:**

   ````bash
   uvicorn main:app --reload --host 0.00port 80   ```

   The API will be available at `http://localhost:8080`
   - API documentation: `http://localhost:8080/docs`
   - Alternative docs: `http://localhost:8080/redoc`
   ````

### Frontend Setup1 **Navigate to the frontend directory:**

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

   The frontend will be available at `http://localhost:30# Available Scripts

### Backend

- `uvicorn main:app --reload` - Start development server with auto-reload

### Frontend

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## API Endpoints

The backend currently provides these endpoints:

- `GET /` - Returns a simple "Hello World message
- `GET /items/{item_id}` - Returns item information with optional query parameter

## Development

### Backend Development

- The backend uses FastAPI for rapid API development
- Auto-generated API documentation is available at `/docs`
- Hot reload is enabled for development

### Frontend Development

- Built with Next.js15nd React 19
- Uses TypeScript for type safety
- Styled with Tailwind CSS v4
- ESLint configured for code quality

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test both backend and frontend
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).
