# client/Dockerfile
FROM nginx:alpine

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy the build folder to Nginx’s default static folder
COPY build/ /usr/share/nginx/html/

# Expose port 80 (Nginx default)
EXPOSE 80

# Start Nginx server (it runs in the foreground by default in this image)
CMD ["nginx", "-g", "daemon off;"]
