# Build stage
FROM node:20-alpine AS build

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and lock files
COPY package.json pnpm-lock.yaml .npmrc ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy the rest of the application
COPY . .

# Build the application
RUN pnpm build

# Runtime stage
FROM node:20-alpine AS runtime

WORKDIR /app

# Install only production dependencies
COPY --from=build /app/package.json /app/pnpm-lock.yaml /app/.npmrc ./
RUN npm install -g pnpm && \
    pnpm install --prod --frozen-lockfile

# Copy build output
COPY --from=build /app/build ./build
COPY --from=build /app/node_modules ./node_modules

# Set default port and allow override via environment variable
ENV PORT=3000
EXPOSE $PORT

# Start the application
CMD ["sh", "-c", "PORT=$PORT node build"] 