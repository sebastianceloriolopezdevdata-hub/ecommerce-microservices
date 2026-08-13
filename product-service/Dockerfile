# ===== Build Stage =====
FROM node:20-alpine AS builder

WORKDIR /app

# Use npm ci for deterministic installations
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# ===== Production Stage =====
FROM node:20-alpine

WORKDIR /app

# Install only production dependencies
COPY package*.json ./
RUN npm ci --omit=dev

# Copy the build output
COPY --from=builder /app/build ./build

# Run directly with node (faster than npm start)
EXPOSE 3000

CMD ["node", "build/index.js"]