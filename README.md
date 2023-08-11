# ClubHandler

## Задание
https://github.com/Mansur51-hub/ClubHandler/blob/84a2a4f9d06156d52ee9fe26a2008a4d87113845/docs/%D0%A2%D0%B5%D1%81%D1%82%D0%BE%D0%B2%D0%BE%D0%B5%20%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5%20GO.pdf

## Запуск решения в docker

<p> В корневую папку репозитория нужно добавить файл для проверки </p>
<p>Далее выполнить команды:</p>

```command
docker build -t {your_image_name} --build-arg file_name={your_file_name} .
```

```command
docker run {your_image_name}
```
