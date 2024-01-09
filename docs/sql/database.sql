drop table if exists user;
create table user (
    id int auto_increment comment 'primarykey',
    name varchar(20) null,
    account varchar(20) null,
    password varchar(255) null,
    invite_code int(8) null comment 'The invitation code in "invite_code" table used when registering',
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    constraint user_pk primary key (id)
) comment 'community user';

drop table if exists invite_code;
create table invite_code (
    id int auto_increment comment 'primarykey',
    code varchar(20) null,
    `state` tinyint(1) null,
    created_at datetime null,
    updated_at datetime null,
    constraint invite_code_pk primary key (id)
) comment 'Cache invite code for register';

