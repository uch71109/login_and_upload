# Requirement

* admin user
  * login feature
  * multiple files upload `local`/`s3`/`gcs`/`blob`
* normal user
  * login feature
  * multiple files upload to local

# System analysis and design

* RPS: `< 10`
* DAU: `2`
* Separate frontend and backend: `Not required`
* Cache server: `Not required`
* Message Queue: `Not required`, but is sutable for files upload to s3/gcs/blob)
* DB: in-memory `sqlite3`

# How to run this project

```
buffalo dev
```

And open url: http://localhost:3000/

## Default user

This project uses in-memory sqlite3 to demo. After run `buffalo dev`, it creates two types of user including `admin user` and `normal user`

|type|email|password|privilege|
|---|---|---|---|
|admin|admin@gmail.com|admin@gmail.com|login/upload videos to `local`/`s3`/`gcs`/`blob`|
|normal|admin@gmail.com|admin@gmail.com|login/upload videos to `local`|
