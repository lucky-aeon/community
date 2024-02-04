alter table types
    add flag_name varchar(10) not null comment '唯一标识名';

alter table types
    add constraint types_pk_flag_name
        unique (flag_name);
