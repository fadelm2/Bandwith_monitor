# Stage 1: Build the React application
FROM node:20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the remaining application source code
COPY . .

# Build the React application for production using Vite
RUN npm run build

# Stage 2: Serve the application using Nginx
FROM nginx:alpine

# Copy the built application from the previous stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Modify the default Nginx configuration to listen on port 8085
RUN sed -i 's/listen  *[0-9]*;/listen 8085;/g' /etc/nginx/conf.d/default.conf

# Expose port 8085
EXPOSE 8085

# Ensure Nginx runs in the foreground
CMD ["nginx", "-g", "daemon off;"]
