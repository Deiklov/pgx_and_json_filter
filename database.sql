create table request_response
(
    guid_transaction uuid default gen_random_uuid() not null,
    guid_strategy uuid default gen_random_uuid() not null,
    filter json,
    request json,
    response json
);

alter table request_response owner to test_user;

create unique index request_response_guid_transaction_guid_strategy_uindex
    on request_response (guid_transaction, guid_strategy);

