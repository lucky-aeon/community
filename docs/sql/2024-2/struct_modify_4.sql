create table article_tags
(
    id          int auto_increment
        primary key,
    tag_name    varchar(20)  null,
    description varchar(255) null,
    user_id     int          not null comment '操作人',
    created_at  datetime     null,
    updated_at  datetime     null,
    deleted_at  datetime     null,
    constraint article_tags_tag_name_uindex
        unique (tag_name)
)
    comment '文章标签';

