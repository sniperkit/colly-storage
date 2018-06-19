# MYSQL

## homebrew

```bash
$ brew install mysql
$ brew info mysql
We've installed your MySQL database without a root password. To secure it run:
    mysql_secure_installation

MySQL is configured to only allow connections from localhost by default

To connect run:
    mysql -uroot

To have launchd start mysql now and restart at login:
  brew services start mysql
Or, if you don't want/need a background service you can just run:
  mysql.server start
 ```


```bash
sudo rm /usr/local/mysql/data/*.err
ps -A | grep -m1 mysql | awk '{print $1}' | sudo xargs kill -9
mysql.server start
```