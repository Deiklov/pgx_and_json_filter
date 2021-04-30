# Nats and docker-compose
Тестирование работы nats c job и cronjob
### Задание  
1. Поднять через docker-compose сервисы: nats, job(3 экземпляра), cronjob
2. Cronjob пишет данные в топик, job читают и выполняют их
3. Смоделировать работу, используя простые реализации job и cronjob
4. Проверить, что сообщения из топика можно забирать в порядке очереди
5. Проверить возможность уведомления jobов о закрытии топика
6. Проверить возможность broadcast сообщения
7. Коннекты должны сбрасываться, когда заканичиваются таски
### Отчет  
1. Запустил 1 cronjob, 1 nats и 3 job (docker-co)
2. pass
3. pass
4. Чтобы jobам не брать одни и те же таски, используем nc.QueueSubscribe,
 тогда они равномерно будут распределяться по jobам без дублирования
5. cronjob кидает в топик сообщение, вида "{flag:fin}", на подписчике вызываем nc.Drain()
6. Сообщения по дефолту идут всем кто подписан на топик,
 можно подписаться на несколько топиков так nc.Subscribe("topic.>") несколько токенов
  или  nc.Subscribe("topic.*") один токен, название топика - (alphanumberic . * >)
7. Сброс коннектов построен по схеме: каждый job отвечает через msg.Respond() to cronjob, cronjob считает кол-во ответов,
 когда соберутся все ответы, то в отдельный топик(не Queue) отправляется сообщение о завершении работы, все подписчики 
 принимают его и сами вызывают nc.Drain()  
 
### Схема работы подсистемы
1. Запускаем ночью cronjob, он выбирает в api ОКР id всех сделок за период, например день
2. По какому нибудь сигналу запускаем несколько jobов
3. Подписиываем всех job через групповую подписку на топик "applicant" и обычную подписку на "notifications"
4. Кидаем через cronjob в топик "requests" сообщения со всеми id сделок
5. Каждый job после обработки сделки отвечает cronjob через reply
6. После того как все сделки обработали, со стороны cronjob в топик "notifications" отправялем сообщения "fin"
7. Закрываем коннект на cronjob и на всех job
8. Завершаем все jobы, ресурсы освобождаются
 		//todo может быть гонка(нету тк таски идут последовательно)
 		//todo ответ если не смогли обработать сообщение
 		//todo проверить когда топик удаляется
 		//todo  добавить запись в postgresql и проверить
 		//todo схема плоха тем что будет бесконечный цикл при ошибке
 		//todo таски вырубаются по таймаутам, крон берет из базы только id
 		
 ### инфа по nats-streaming 
 1. нет request-reply вместо этого модифицированный ack  
 2. тк через ack не можем скинуть доп инфу, только ок или !ок, то непонтяно почему !ок и как это обработать
 3. maxinflight покажет сколько сообщений можно закинуть неподтвержденными к 1 подписчику  
 4. topic nats and nats streaming не пересекаются, те subcribe("appliciant") будут принимать разные сообщения  
 
 
### Схема работы v2  
1. Джобы так же подпсываются через QueueSubscribe()
2. Если нет новых задач в течении 10 min, то джоба заврешается
(соответсвенно все его горутины тоже свернуться, поэтому 10 мин отсчитываем от взятие последнего таска)  
3. Добавляем stan.SetManualAckMode(), stan.AckWait(30*time.Second) для возврата в очередь плохо обработанных задач
4. Cron должен узнать что таск отработал, через msg.Respond() уже не получится, можно через openshift,  
как увидим что ок, то вычеркиваем джобы из базы, но не понятно какие таски обработались(нужен ответ от каждой горутины)  
