# Energy War Game

A "battleship" like game but with power plants. Each user has to setup its power plants (nuclear, gas, wind or solar) to meet the required capacity. Players take turns striking to reduce their opponent's capacity. If a plant is hit, all the capacity is removed.

## Docker Deployment

### Building and Running Locally

1. Build the Docker image:
   ```bash
   docker build -t energywar .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 energywar
   ```

   Or using Docker Compose:
   ```bash
   docker-compose up
   ```

3. Access the application at http://localhost:8080

### Deploying to Digital Ocean

#### Option 1: Using Docker Hub

1. Build and tag the Docker image:
   ```bash
   docker build -t yourusername/energywar:latest .
   ```

2. Push the image to Docker Hub:
   ```bash
   docker login
   docker push yourusername/energywar:latest
   ```

3. Create a Droplet on Digital Ocean:
   - Choose the Docker image from the marketplace
   - Select your preferred size and region
   - Add your SSH key

4. SSH into your Droplet:
   ```bash
   ssh root@your-droplet-ip
   ```

5. Pull and run the Docker image:
   ```bash
   docker pull yourusername/energywar:latest
   docker run -d -p 80:8080 yourusername/energywar:latest
   ```

#### Option 2: Using Digital Ocean App Platform

1. Create a repository on GitHub or GitLab and push your code.

2. Log in to Digital Ocean and navigate to the App Platform.

3. Click "Create App" and select your repository.

4. Configure the app:
   - Select the branch to deploy
   - Choose the Dockerfile as the build method
   - Configure the HTTP port to 8080
   - Set up any environment variables if needed

5. Choose your plan and click "Launch App".

#### Option 3: Using Digital Ocean Container Registry

1. Create a Container Registry on Digital Ocean.

2. Install and configure the `doctl` CLI:
   ```bash
   doctl auth init
   ```

3. Log in to the Digital Ocean Container Registry:
   ```bash
   doctl registry login
   ```

4. Build and tag the Docker image:
   ```bash
   docker build -t registry.digitalocean.com/your-registry/energywar:latest .
   ```

5. Push the image to the Digital Ocean Container Registry:
   ```bash
   docker push registry.digitalocean.com/your-registry/energywar:latest
   ```

6. Deploy to a Droplet or Kubernetes cluster on Digital Ocean.

## Game Rules

| Power plant | Code | Capacity | size  |
| ----------- | ---- | -------- | ----- |
| NUCLEAR     | N    | 1000     | 3 x 3 |
| GAS         | G    | 300      | 2 x 2 |
| WIND        | W    | 100      | 2 x 1 |
| SOLAR       | S    | 25       | 1 x 1 |

### Mechanics
- User should build an energy infrastructure that meets at least the capacity defined in the game and max a 10% extra of the capacity
- If a power plant is HIT, capacity of the entire plant is removed from the counter
- The game ends when one of the players have below the 10% of the defined capacity

## API Documentation

The API documentation is available at `/swagger/index.html` when the application is running.
