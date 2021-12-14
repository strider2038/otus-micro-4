# Otus Microservice architecture Homework 5

## Домашнее задание выполнено для курса ["Microservice architecture"](https://otus.ru/lessons/microservice-architecture/)

### Запуск приложения

```shell
# запуск minikube
# версия k8s v1.19, на более поздних есть проблемы с установкой ambassador
minikube start --cpus=6 --memory=8g --vm-driver=virtualbox --cni=flannel --kubernetes-version="v1.19.0"

kubectl create namespace otus

# установка ambassador
helm install aes datawire/ambassador -f deploy/ambassador-values.yaml

## запуск проекта
helm install --wait -f deploy/identity-values.yaml identity-service ./services/identity-service/deployments/identity-service --atomic
helm install --wait echo-service ./services/echo-service/deployments/echo-service --atomic

# применение настроек ambassador
kubectl apply -f services/ambassador/
```

### Тестирование

Тесты Postman расположены в директории `test/postman`. Запуск тестов.

```bash
newman run ./test/postman/test.postman_collection.json
```

Или с использованием Docker.

```bash
docker run -v $PWD/test/postman/:/etc/newman --network host -t postman/newman:alpine run test.postman_collection.json
```

## Диаграммы взаимодействия сервисов

### Регистрация

![Регистрация](docs/images/1-registration.png)

### Аутентификация

![Аутентификация](docs/images/2-authentication.png)

### Доступ к профилю

![Доступ к профилю](docs/images/3-get-profile.png)

### Запрос в backend сервис

![Запрос в backend сервис](docs/images/4-access-backend.png)
