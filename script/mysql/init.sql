create database if not exists `integration_test`;
create table if not exists `integration_test`.`test_model`(
    `id` bigint auto_increment,
    `first_name` varchar(1024) not null,
    `age` smallint not null,
    `last_name` varchar(1024),
    primary key (`id`)
);
