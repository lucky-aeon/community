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

drop table if exists article;
create table article (
    id int auto_increment comment 'primarykey',
    title varchar(50) null,
    `description` varchar(5000) null,
    user_id int not null,
    issue_id int default 0,
    solved tinyint(1) default 0,
    `like` int default 0,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    constraint article_pk primary key (id)
) comment 'it is issue or answer';