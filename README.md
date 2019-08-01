# Deployer

Deploy applications using [GitHub Deployments API](https://developer.github.com/v3/repos/deployments/).

1. Listen/poll new GitHub deployments
2. Create pending GitHub deployment status and empty Gist
3. Execute command file that deploys application
4. Create success/failure Github deployment status and update Gist
5. Sleep 30 seconds
6. Go to step #1

To register github deployments you can use curl,
[slashdeploy](https://github.com/remind101/slashdeploy) or something else.

If you are using Slack, you can enable 'Deploy Events -> Show deployment statuses'
to get slack notifications with links to Gist.

## Usage

```bash
deployer -h
```

```
  -env string
      Github deployment environment (required)
  -file string
      File to execute when new deployment is available (required)
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

## Example

Trigger deployment from developer's machine:

```bash
$ deploy --branch v49.3 --env demo
```

On the application server:

```bash
$ deployer -u ghuser -p ghpass -org mak-it -repo myapp -env demo -file bin/deploy-stub -sleep 30
```

Output:

```
Deployer
ghuser ghpass mak-it myapp demo myapp-demo
Sleep duration: 30s
2017/08/24 13:32:52 Getting deployment  demo
2017/08/24 13:32:53 Problem in getting deployment Deployment statuses already present
2017/08/24 13:32:53 Sleeping 30s
2017/08/24 13:33:23 Getting deployment  demo
2017/08/24 13:33:23 123 v49.3 demo
2017/08/24 13:33:24 Gist create ID: 123, Ref: v49.3, Environment: demo
2017/08/24 13:33:25 Gist https://gist.github.com/...
2017/08/24 13:33:25 Deployment status create  pending
2017/08/24 13:33:26 Executing File bin/deploy-stub
-----> Deploying version: v49.3...
-----> Adding BUILD_ENV to build environment...
-----> Compiling Ruby/Rails
=====> Application deployed:
       http://ruby-rails-sample.dokku.me
2017/08/24 13:33:26 Deployment status create  success
2017/08/24 13:33:26 Gist update
```

## Build

```bash
make
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/mak-it/deployer. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

Deployer is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
