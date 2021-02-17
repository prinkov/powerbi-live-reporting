For run

```
docker build -t myapp:latest . &&  docker run --env-file .env --name pbi_publisher myapp:latest 
```