
alter table message_states
    add state tinyint(1) DEFAULT '0';
    add type tinyint(1) DEFAULT '1';
    add articleId int(11) NOT NULL;


alter table message_logs
    add articleId int(11) NOT NULL;