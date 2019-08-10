

# Feature Description
1. run as web server
2. save tasks in redis, then when server panic and restart, tasks will restore.
3. all added tasks organized by app specified key
4. add tasks through http calling
5. tasks fired by calling url with parameters
6. tasks fired result save in log file

# Tasks Type
1. delayed tasks fired once
2. delayed tasks fired repeated

