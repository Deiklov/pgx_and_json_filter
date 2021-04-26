# Nats and docker-compsoe
Тестирование работы nats c job и cronjob
### Задание  
1. Поднять через docker-compose сервисы: nats, job(3 экземпляра), cronjob
2. Cronjob пишет данные в топик, job читают и выполняют их
3. Смоделировать работу, используя простые реализации jon и cronjob
4. Проверить, что сообщения из топика можно забирать в порядке очереди
5. Проверить возможность уведомления jobов о закрытии топика
6. Проверить возможность broadcast сообщения
7. Коннекты должны сбрасываться, когда заканичиваются таски
### Отчет  
1. 
4. Чтобы jobам не брать одни и те же таски, используем nc.QueueSubscribe,
 тогда они равномерно будут распределяться по jobам без дублирования
5. cronjob кидает в топик сообщение, вида "{flag:fin}", на подписчике вызываем nc.Drain()
6. Сообщения по дефолту идут всем кто подписан на топик,
 можно подписаться на несколько топиков так nc.Subscribe("topic.>") или  nc.Subscribe("topic.*"), название топика - (alphanumberic . * >)
 