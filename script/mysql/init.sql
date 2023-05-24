create database if not exists `integration_test`;
create table if not exists `integration_test`.`test_model`(
    `id` bigint auto_increment,
    `first_name` varchar(1024) not null,
    `age` smallint not null,
    `last_name` varchar(1024),
    primary key (`id`)
);
CREATE TABLE IF NOT EXISTS `integration_test`.`order`(
    `id` bigint auto_increment,
    `using_col1` varchar(1024) not null,
    `using_col2` varchar(1024) not null,
    primary key (`id`)
);

CREATE TABLE IF NOT EXISTS `integration_test`.`order_detail`(
    `order_id` bigint not null,
    `item_id` bigint not null,
    `using_col1` varchar(1024) not null,
    `using_col2` varchar(1024) not null
);

CREATE TABLE IF NOT EXISTS `integration_test`.`item`(
    `id` bigint auto_increment,
    primary key (`id`)
);