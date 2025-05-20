# Build stage
FROM node:20-alpine AS builder
WORKDIR /app

# Define build arguments
ARG SERVICE_API_HOST
ARG VITE_SERVICE_API_HOST
ARG PORT

# Set build-time environment variables from ARGs
ENV SERVICE_API_HOST=$SERVICE_API_HOST
ENV VITE_SERVICE_API_HOST=$VITE_SERVICE_API_HOST
ENV PORT=$PORT

# Install pnpm
RUN npm install -g pnpm

# Install dependencies
COPY package.json ./
RUN pnpm install

# Copy source code
COPY . .

# Build the SvelteKit app
RUN pnpm build

# Production stage
FROM node:20-alpine
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy necessary files from the builder stage
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Expose the port the app runs on
EXPOSE ${PORT}

# Set runtime environment variables
ENV NODE_ENV=production
ENV SERVICE_API_HOST=$SERVICE_API_HOST
ENV VITE_SERVICE_API_HOST=$VITE_SERVICE_API_HOST
ENV PORT=$PORT

# Start the application
CMD ["node", "build"]
