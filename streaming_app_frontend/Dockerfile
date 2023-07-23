# Use the official Node.js image as the base image
# Build stage
FROM node:18.12.1-alpine AS build-stage

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json into the container
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code into the container
COPY . .

# Build the application
RUN npm run build

# Print the contents of the 'dist' directory
RUN ls /app/dist

# Production stage
FROM nginx:1.12-alpine

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Create a new user and group called "app"
RUN addgroup app && adduser -S -G app app

# Change the ownership and permissions of the nginx cache directory
RUN chown -R app:app /var/cache/nginx

# Create the PID file and grant permission to the app user and group
RUN touch /var/run/nginx.pid && \
    chown -R app:app /var/run/nginx.pid

# Copy the Nginx configuration
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

USER app

# Copy the built application code from the build stage
COPY --from=build-stage /app/dist /usr/share/nginx/html


# Expose the port that the container will listen on
EXPOSE 3000

# Start nginx and keep the container running
CMD ["nginx", "-g", "daemon off;"]
