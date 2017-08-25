# Deployer

Deploy [dokku](http://dokku.viewdocs.io/dokku/) applications using [GitHub Deployments API](https://developer.github.com/v3/repos/deployments/).

1. Listen/poll new GitHub deployments
2. Create pending GitHub deployment status and empty Gist
3. Deploy application to dokku
4. Create success/failure Github deployment status and update Gist
5. Sleep 30 seconds
6. Go to step #1

Make sure to run `deployer` as `dokku` user.

To register github deployments you can use [deploy](https://github.com/remind101/deploy),
[slashdeploy](https://github.com/remind101/slashdeploy) or something else.

If you are using Slack, you can enable 'Deploy Events -> Show deployment statuses'
to get slack notifications with links to Gist.

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

## Installation

Install `deployer` on dokku server:

```bash
$ su - dokku
$ curl -L -o deployer https://github.com/mak-it/deployer/releases/download/v0.1.1/deployer-linux-amd64
$ chmod +x deployer
```

Start `deployer` on system reboot using `cron` and `screen`:

```bash
$ su - dokku
$ echo "@reboot /usr/bin/screen -d -m /home/dokku/deployer -u ghuser -p ghpass -org mak-it -repo myapp -env demo -app myapp-demo -sleep 30" | crontab -
```

Start `deployer` in `screen`:

```bash
$ su - dokku
$ /usr/bin/screen -d -m /home/dokku/deployer -u ghuser -p ghpass -org mak-it -repo myapp -env demo -app myapp-demo -sleep 30
```

## Example

Trigger deployment from developer's machine:

```bash
$ cd ~/src/myapp
$ deploy --branch v49.3 --env demo
```

On the dokku server:

```bash
$ deployer -u ghuser -p ghpass -org mak-it -repo myapp -env demo -app myapp-demo -sleep 30
```

Output:

```
Deployer
ghuser ghpass mak-it myapp demo myapp-demo
Sleep duration: 30s
2017/08/24 13:31:43 Cloning repo  
2017/08/24 13:31:51 Adding repo dokku remote  .myapp-demo myapp-demo
2017/08/24 13:32:52 Getting deployment  demo
2017/08/24 13:32:53 Problem in getting deployment Deployment statuses already present
2017/08/24 13:32:53 Sleeping 30s
2017/08/24 13:33:23 Getting deployment  demo
2017/08/24 13:33:23 123 v49.3 demo
2017/08/24 13:33:23 Fetchin repo  .myapp-demo
2017/08/24 13:33:24 Gist create ID: 123, Ref: v49.3, Environment: demo
2017/08/24 13:33:25 Gist https://gist.github.com/...
2017/08/24 13:33:25 Deployment status create  pending
2017/08/24 13:33:26 Deploying repo  .myapp-demo
Everything up-to-date
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
