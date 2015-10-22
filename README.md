# bane

[![Circle CI](https://circleci.com/gh/jfrazelle/bane.svg?style=svg)](https://circleci.com/gh/jfrazelle/bane)

AppArmor profile generator for docker containers. Basically a better AppArmor
profile, than creating one by hand, because who would ever do that.

> "Reviewing AppArmor profile pull requests is the _bane_ of my existence"
>  Jess Frazelle


```console
$ bane -h
 _
| |__   __ _ _ __   ___
| '_ \ / _` | '_ \ / _ \
| |_) | (_| | | | |  __/
|_.__/ \__,_|_| |_|\___|
 Custom AppArmor profile generator for docker containers
 Version: v0.1.0

  -d    run in debug mode
  -profile-dir string
        directory for saving the profiles (default "/etc/apparmor.d/containers")
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

### Config File

Below is the sample config for nginx in a container:

```toml
# name of the profile, we will auto prefix with `docker-`
# so the final profile name will be `docker-nginx`
Name = "nginx"

[Filesystem]
# read only paths for the container
ReadOnlyPaths = [
	"/bin/**",
	"/bin/**",
	"/boot/**",
	"/dev/**",
	"/etc/**",
	"/home/**",
	"/lib/**",
	"/lib64/**",
	"/media/**",
	"/mnt/**",
	"/opt/**",
	"/proc/**",
	"/root/**",
	"/sbin/**",
	"/srv/**",
	"/tmp/**",
	"/sys/**",
	"/usr/**",
]

# paths where you want to log on write
LogOnWritePaths = [
	"/**"
]

# paths where you can write
WritablePaths = [
	"/var/run/nginx.pid"
]

# allowed executable files for the container
AllowExec = [
	"/usr/sbin/nginx"
]

# denied executable files
DenyExec = [
	"/bin/dash",
	"/bin/sh",
	"/usr/bin/top"
]

[Network]
# if you don't need to ping in a container, you can probably
# set Raw to false and deny network raw
Raw = false
Packet = false
```

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
$ bane sample.toml
# Profile installed successfully you can now run the profile with
# `docker run --security-opt="apparmor:docker-nginx"`

# now let's run nginx
$ docker -d run --security-opt="apparmor:docker-nginx" -p 80:80 nginx
```

Using custom AppArmor profiles has never been easier!

**Now let's try to do malicious activites with the sample profile:**

```console
$ docker run --security-opt="apparmor:docker-nginx" -p 80:80 --rm -it nginx bash
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

For the above `sample.toml` the generated profile is:

```
#include <tunables/global>

profile docker-nginx flags=(attach_disconnected,mediate_deleted) {
  #include <abstractions/base>

  network,
  deny network raw,
  deny network packet,

  capability,
  file,
  umount,

  deny /bin/** wl,
  deny /bin/** wl,
  deny /boot/** wl,
  deny /dev/** wl,
  deny /etc/** wl,
  deny /home/** wl,
  deny /lib/** wl,
  deny /lib64/** wl,
  deny /media/** wl,
  deny /mnt/** wl,
  deny /opt/** wl,
  deny /proc/** wl,
  deny /root/** wl,
  deny /sbin/** wl,
  deny /srv/** wl,
  deny /tmp/** wl,
  deny /sys/** wl,
  deny /usr/** wl,

  audit /** w,

  /var/run/nginx.pid w,

  /usr/sbin/nginx ix,

  deny /bin/dash mrwklx,
  deny /bin/sh mrwklx,
  deny /usr/bin/top mrwklx,

  deny @{PROC}/{*,**^[0-9*],sys/kernel/shm*} wkx,
  deny @{PROC}/sysrq-trigger rwklx,
  deny @{PROC}/mem rwklx,
  deny @{PROC}/kmem rwklx,
  deny @{PROC}/kcore rwklx,
  deny mount,
  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/efi/efivars/** rwklx,
  deny /sys/kernel/security/** rwklx,
}
```

## TODO

- add all the network controls like `tcp` etc
- more tunables
- add capabilities
- add syscalls
- tests (integration, unit)
