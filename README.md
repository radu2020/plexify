# plexify
Job Processing Server

# Running the server and simulation
## Instalation
```
git clone https://github.com/radu2020/plexify.git
cd plefixy
```

## Build
```
make build
```

## Run Server and Client
Run in separate windows side-by-side

```
make server
make client
```

## Clean
```
make clean
```

# Assumptions and trade-offs.
Worker Pool Size: Fixed in the code but can be adjusted based on the expected load.

Channel Buffer: A buffered channel of size 100 is used for the job queue.

Random Processing Time: Simulated with time.Sleep for simplicity. In production, actual processing logic would replace this.

Status Tracking: Uses sync.Map for simplicity. For large-scale systems, a database would be more appropriate.