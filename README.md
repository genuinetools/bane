# bane

[![Travis CI](https://travis-ci.org/jessfraz/bane.svg?branch=master)](https://travis-ci.org/jessfraz/bane)

AppArmor profile generator for docker containers. Basically a better AppArmor
profile, than creating one by hand, because who would ever do that.

> "Reviewing AppArmor profile pull requests is the _bane_ of my existence"
>  - Jess Frazelle

![bane](bane.jpg)


## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-darwin-386) / [amd64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-freebsd-386) / [amd64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-linux-386) / [amd64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-linux-amd64) / [arm](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-linux-arm) / [arm64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-windows-386) / [amd64](https://github.com/jessfraz/bane/releases/download/v0.2.2/bane-windows-amd64)

#### Via Go

```bash
$ go get github.com/jessfraz/bane
```

## Usage

```console
$ bane -h
 _
| |__   __ _ _ __   ___
| '_ \ / _` | '_ \ / _ \
| |_) | (_| | | | |  __/
|_.__/ \__,_|_| |_|\___|
 Custom AppArmor profile generator for docker containers
 Version: v0.2.2

  -d    run in debug mode
  -profile-dir string
        directory for saving the profiles (default "/etc/apparmor.d/containers")
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

### Config File

[sample.toml](sample.toml) is a AppArmor sample config for nginx in a container.

#### File Globbing

| Glob Example  | Description |
| ------------- | ------------- |
| `/dir/file` |   match a specific file |
| `/dir/*`        | match any files in a directory (including dot files) |
| `/dir/a*`      | match any file in a directory starting with a |
| `/dir/*.png`    | match any file in a directory ending with .png |
| `/dir/[^.]*`   | match any file in a directory except dot files |
| `/dir/`        | match a directory |
| `/dir/*/`       | match any directory within /dir/ |
| `/dir/a*/`     | match any directory within /dir/ starting with a |
| `/dir/*a/`     | match any directory within /dir/ ending with a |
| `/dir/**`       | match any file or directory in or below /dir/ |
| `/dir/**/`     | match any directory in or below /dir/ |
| `/dir/**[^/]`   | match any file in or below /dir/ |
| `/dir{,1,2}/**` | match any file or directory in or below /dir/, /dir1/, and /dir2/ |

### Installing a Profile

Now that we have our config file from above let's install it. `bane` will
automatically install the profile in a directory
`/etc/apparmor.d/containers/` and run `apparmor_parser`.

```console
$ sudo bane sample.toml
# Profile installed successfully you can now run the profile with
# `docker run --security-opt="apparmor:docker-nginx-sample"`

# now let's run nginx
$ docker run -d --security-opt="apparmor:docker-nginx-sample" -p 80:80 nginx
```

Using custom AppArmor profiles has never been easier!

**Now let's try to do malicious activities with the sample profile:**

```console
$ docker run --security-opt="apparmor:docker-nginx-sample" -p 80:80 --rm -it nginx bash
root@6da5a2a930b9:~# ping 8.8.8.8
ping: Lacking privilege for raw socket.

root@6da5a2a930b9:/# top
bash: /usr/bin/top: Permission denied

root@6da5a2a930b9:~# touch ~/thing
touch: cannot touch 'thing': Permission denied

root@6da5a2a930b9:/# sh
bash: /bin/sh: Permission denied

root@6da5a2a930b9:/# dash
bash: /bin/dash: Permission denied
```


Sample `dmesg` output when using `LogOnWritePaths`:

```
[ 1964.142128] type=1400 audit(1444369315.090:38): apparmor="STATUS" operation="profile_replace" profile="unconfined" name="docker-nginx" pid=3945 comm="apparmor_parser"
[ 1966.620327] type=1400 audit(1444369317.570:39): apparmor="AUDIT" operation="open" profile="docker-nginx" name="/1" pid=3985 comm="nginx" requested_mask="c" fsuid=0 ouid=0
[ 1966.624381] type=1400 audit(1444369317.574:40): apparmor="AUDIT" operation="mkdir" profile="docker-nginx" name="/var/cache/nginx/client_temp/" pid=3985 comm="nginx" requested_mask="c" fsuid=0 ouid=0
[ 1966.624446] type=1400 audit(1444369317.574:41): apparmor="AUDIT" operation="chown" profile="docker-nginx" name="/var/cache/nginx/client_temp/" pid=3985 comm="nginx" requested_mask="w" fsuid=0 ouid=0
[ 1966.624463] type=1400 audit(1444369317.574:42): apparmor="AUDIT" operation="mkdir" profile="docker-nginx" name="/var/cache/nginx/proxy_temp/" pid=3985 comm="nginx" requested_mask="c" fsuid=0 ouid=0
[ 1966.624494] type=1400 audit(1444369317.574:43): apparmor="AUDIT" operation="chown" profile="docker-nginx" name="/var/cache/nginx/proxy_temp/" pid=3985 comm="nginx" requested_mask="w" fsuid=0 ouid=0
[ 1966.624507] type=1400 audit(1444369317.574:44): apparmor="AUDIT" operation="mkdir" profile="docker-nginx" name="/var/cache/nginx/fastcgi_temp/" pid=3985 comm="nginx" requested_mask="c" fsuid=0 ouid=0
[ 1966.624534] type=1400 audit(1444369317.574:45): apparmor="AUDIT" operation="chown" profile="docker-nginx" name="/var/cache/nginx/fastcgi_temp/" pid=3985 comm="nginx" requested_mask="w" fsuid=0 ouid=0
[ 1966.624546] type=1400 audit(1444369317.574:46): apparmor="AUDIT" operation="mkdir" profile="docker-nginx" name="/var/cache/nginx/uwsgi_temp/" pid=3985 comm="nginx" requested_mask="c" fsuid=0 ouid=0
[ 1966.624582] type=1400 audit(1444369317.574:47): apparmor="AUDIT" operation="chown" profile="docker-nginx" name="/var/cache/nginx/uwsgi_temp/" pid=3985 comm="nginx" requested_mask="w" fsuid=0 ouid=0
```


### What does the generated profile look like?

For the above `sample.toml` the generated profile is available as [docker-nginx-sample](docker-nginx-sample).

### Integration with Docker

This was originally a proof of concept for what will hopefully become a native
security profile in the Docker engine. For more information on this, see
[docker/docker#17142](https://github.com/docker/docker/issues/17142).

