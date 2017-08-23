# Deployer

###

1. Listen to new GitHub deployments
2. Create pending GitHub deployment status and empty Gist
3. Deploy application to dokku and update Gist
4. Create success/failure Github deployment status
5. Sleep
6. Go to step #1

## Usage

```bash
deployer -h
```

```
  -app string
    	Dokku application name (required)
  -env string
    	Github deployment environment (required)
  -org string
    	GitHub org (required)
  -p string
    	GitHub password (required)
  -repo string
    	GitHub repo (required)
  -sleep int
    	Time to sleep between loops (defaults to 30 seconds) (default 30)
  -u string
    	GitHub username (required)
```

## Build

```bash
make
```
