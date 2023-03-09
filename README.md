![intro](https://user-images.githubusercontent.com/126064091/224017711-21febf1f-83d8-4ff9-a77d-bcfb2f660824.svg)

# About this project:
This project is a pet project of mine that I decided to do for backend studies.

It is a REST API that generates a translation of an English word into Russian using Yandex Dictionary and generates an image based on that word using OpenAI's DALL-E model.

## Design:
![design drawio](https://user-images.githubusercontent.com/126064091/224037889-37add257-2373-4049-8310-63a07f3f9d41.png)

## API Documentation:

## Limitations:
Creating pictures using DALL-E is limited to 50 requests per minute, so there is a limit inside the code.

# How to use it:
Build all containers using-docker compose
```bash
make build 
```
Run the API perfomance test with K6
```bash
make stress
```
After starting the test, the following Grafana GUI opens in the browser:
![image](https://user-images.githubusercontent.com/126064091/224057841-17c02fac-6b64-4451-bb75-91ed32b7ff16.png)

The number of virtual users and the number of requests per second can be configured in the docker-compose config from the repository:
```yaml
 k6:
  image: grafana/k6:latest
  profiles: ["k6"]
  command: run --vus 100 --rps 500 --duration 30s /test.js
  volumes:
   - ./test.js:/test.js
```
