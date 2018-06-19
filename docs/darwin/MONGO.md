# MONGO

## homebrew

```bash
$ brew install mongodb
$ brew info mongodb
mongodb: stable 3.6.5 (bottled)
High-performance, schema-free, document-oriented database
https://www.mongodb.org/
/usr/local/Cellar/mongodb/3.6.5 (19 files, 297.9MB) *
  Poured from bottle on 2018-06-19 at 13:09:11
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/mongodb.rb
==> Dependencies
Build: go ✔, pkg-config ✔, scons ✘
Required: python@2 ✔
Recommended: openssl ✔
Optional: boost ✔
==> Requirements
Build: xcode ✔
Required: macOS >= 10.8 ✔
==> Options
--with-boost
	Compile using installed boost, not the version shipped with mongodb
--with-sasl
	Compile with SASL support
--without-openssl
	Build without openssl support
==> Caveats
To have launchd start mongodb now and restart at login:
  brew services start mongodb
Or, if you don't want/need a background service you can just run:
  mongod --config /usr/local/etc/mongod.conf
```