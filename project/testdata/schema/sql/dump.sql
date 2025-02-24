create database if not EXISTS appdb default character set utf8;
USE appdb;

CREATE TABLE IF NOT EXISTS users (  
    id bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Уникальный идентификатор пользователя',  
    last_name VARCHAR(255) NOT NULL COMMENT 'Фамилия пользователя',  
    first_name VARCHAR(255) NOT NULL COMMENT 'Имя пользователя',  
    father_name VARCHAR(255) COMMENT 'Отчество пользователя (может быть NULL, если отсутствует)',  
    years_old INT NOT NULL COMMENT 'Возраст пользователя',  
    gender_code VARCHAR(16) NOT NULL COMMENT 'Пол пользователя: man - мужчина, woman - женщина',  
    bank VARCHAR(255) NOT NULL COMMENT 'Название банка, связанного с пользователем',
    PRIMARY KEY (id)
) COMMENT='Таблица для хранения информации о пользователях'; 

CREATE TABLE IF NOT EXISTS events (  
    id bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Уникальный идентификатор события',  
    price INT NOT NULL COMMENT 'Цена услуги, связанной с событием',  
    service_id INT NOT NULL COMMENT 'Идентификатор услуги, связанной с событием',  
    location_id INT NOT NULL COMMENT 'Идентификатор локации, где произошло событие',  
    timestamp DATETIME NOT NULL COMMENT 'Дата и время события',
    PRIMARY KEY (id)
) COMMENT='Таблица для хранения информации о событиях';  