# REDIS

## homebrew

```bash
$ brew install redis
$ brew install hiredis
$ brew info redis
redis: stable 4.0.10 (bottled), HEAD
Persistent key-value database, with built-in net interface
https://redis.io/
/usr/local/Cellar/redis/4.0.10 (13 files, 2.8MB) *
  Poured from bottle on 2018-06-18 at 11:59:15
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/redis.rb
==> Options
--with-jemalloc
	Select jemalloc as memory allocator when building Redis
--HEAD
	Install HEAD version
==> Caveats
To have launchd start redis now and restart at login:
  brew services start redis
Or, if you don't want/need a background service you can just run:
  redis-server /usr/local/etc/redis.conf
```