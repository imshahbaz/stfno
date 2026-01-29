# Use an official lightweight Python image
FROM python:3.11-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the requirements file first to leverage Docker cache
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of your application code
COPY . .

# Expose the port FastAPI will run on
EXPOSE 8090

# Command to run the application using Uvicorn
# We bind to 0.0.0.0 so Northflank can route external traffic to the container
CMD ["uvicorn", "main:app", "--port", "8090"]