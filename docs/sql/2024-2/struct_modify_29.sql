alter table users
    add subscribe tinyint(1) NOT NULL DEFAULT 0;

ALTER TABLE subscriptions CHANGE `desc` content int(11) ;