# Base stage for both development and production
FROM node:20-alpine AS base
WORKDIR /app

# Dependencies stage - used by both dev and prod
FROM base AS dependencies
COPY package*.json ./
RUN npm ci

# Development stage
FROM base AS development
COPY --from=dependencies /app/node_modules ./node_modules
# Copy only package files first to optimize rebuilds
COPY package*.json ./
# We'll mount the rest of the code as a volume in docker-compose

# Expose port
EXPOSE 3000

# Set environment variables for development
ENV NODE_ENV=development
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Run development server with hot reloading
CMD ["npm", "run", "dev"]

# Build stage for production
FROM base AS builder
COPY --from=dependencies /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Production stage
FROM base AS production
COPY --from=dependencies /app/node_modules ./node_modules
# Only copy build artifacts and necessary files for production
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./

# Set environment variables for production
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Run production server
CMD ["npm", "run", "start"]