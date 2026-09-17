create table person
(
    person_id varchar(64)             not null
        primary key,
    login     varchar(64)             not null,
    full_name varchar(255) default '' not null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

create table session
(
    session_id  varchar(64)                              not null
        primary key,
    person_id   varchar(64)                              not null,
    created_dt  datetime     default current_timestamp() null,
    expiry_dt   datetime                                 not null,
    expired_ind varchar(1)   default 'N'                 null,
    token       varchar(256) default ''                  not null,
    constraint fk_session_user
        foreign key (person_id) references person (person_id)
)CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

create table urls
(
    url_id         varchar(64)   not null
        primary key,
    person_id      varchar(64)   not null,
    shortened_code varchar(64)   not null,
    full_url       varchar(1024) not null,
    constraint fk_urls_person
        foreign key (person_id) references person (person_id)
)CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

