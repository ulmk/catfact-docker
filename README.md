## To Run Docker

### Build image

```
docker build -t go-app .

```

### Run image
```
docker run --name=go-web-app -p 80:5500 -t go-app

```
