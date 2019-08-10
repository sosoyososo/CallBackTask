

# Feature Description
1. run as web server
2. save tasks in mysql, then when server panic and restart, tasks will restore.
3. all added tasks organized by app specified key
4. add tasks through http calling
5. tasks fired by calling url with parameters
6. tasks fired result save in log file

# Tasks Type
1. delayed tasks fired once
2. delayed tasks fired repeated

# Services 
```
// list tasks with group key default
curl http://localhost:8080/listTask?group=default 

// add single fired task
curl -d '{"groupKey":"default","delay":5,"duration":2,"repeat":false,"callBackURL":"localhost"}' http://localhost:8080/addTask 

// add cycle fired task
curl -d '{"groupKey":"default","delay":5,"duration":2,"repeat":true,"callBackURL":"localhost"}' http://localhost:8080/addTask 

//cancel task 
curl http://localhost:8080/cancelTask?id=812f5214-ad6a-4d18-ac74-b1c904868cb8

```

# Config Example
```
// debug mode config file
// conf-dev.json 
{
  "mysqlHost":"127.0.0.1",
  "mysqlUserName":"root",
  "mysqlUserPswd":"root",
  "serverPort":":8081"
}

// debug mode config file
// conf.json 
{
  "mysqlHost":"127.0.0.1",
  "mysqlUserName":"root",
  "mysqlUserPswd":"root",
  "serverPort":":8081"
}

```