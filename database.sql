create table request_response
(
    guid_transaction uuid default gen_random_uuid() not null,
    guid_strategy    uuid default gen_random_uuid() not null,
    filter           jsonb,
    request          jsonb,
    response         jsonb
);

alter table request_response
    owner to test_user;

create unique index request_response_guid_transaction_guid_strategy_uindex
    on request_response (guid_transaction, guid_strategy);
-- https://postgrespro.ru/docs/postgrespro/9.5/datatype-json
-- используем gin индекс на основе оф доки
-- CREATE INDEX idxgintags ON api USING GIN ((jdoc -> 'tags'));
-- есть параметр fastupdate(on,off), лимит для gin_pending_list_limit(объем памяти для перестроения индекса)

create index idxgin on request_response using gin (request JSONB_PATH_OPS);
-- Пример запросов
select guid_transaction,
       guid_strategy,
       request -> 'Applicant' -> 'Cur_City',
       (request -> 'Applicant' -> 'APP_Inner_ConclDate')
from request_response
WHERE (request -> 'Applicant' -> 'Cur_City') not in ('["Янаул"]', '["Якутск"]', '["Южно-Сахалинск"]')
   or (request -> 'Applicant' -> 'APP_Inner_ConclDate') ?| array ['2019-11-11T00:00:00Z','2018-02-10T00:00:00Z']