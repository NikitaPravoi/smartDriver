create table if not exists orders
(
    id             serial
        primary key,
    customer_name  varchar(255)   not null,
    customer_phone varchar(20)    null,
    city           varchar(64)    null,
    street         varchar(64)    null,
    apartment      varchar(64)    null,
    floor          int            null,
    entrance       int            null,
    comment        varchar(255)   null,
    cost           numeric(10, 2) not null,
    status         int default 0  null,
    created_at     timestamp      not null
);

create table if not exists organizations
(
    id      serial
        primary key,
    name    varchar(128)    not null,
    balance real default 0 not null,
    iiko_id varchar(256)    not null,
    constraint organizations_iiko_id_uindex
        unique (iiko_id),
    constraint organizations_pk
        unique (iiko_id)
);

create table if not exists branches
(
    id              serial
        primary key,
    name            varchar(64) not null,
    location        point       not null,
    organization_id int         not null,
    constraint branches_organizations_id_fk
        foreign key (organization_id) references organizations (id)
);

create table if not exists plans
(
    id             serial
        primary key,
    name           varchar(255)   not null,
    cost           numeric(10, 2) not null,
    employee_limit int            not null
);

create table if not exists organization_plans
(
    id              serial
        primary key,
    organization_id int  not null,
    plan_id         int  not null,
    start_date      date not null,
    end_date        date null,
    constraint organization_plans_ibfk_1
        foreign key (organization_id) references organizations (id),
    constraint organization_plans_ibfk_2
        foreign key (plan_id) references plans (id)
);

create index organization_id_idx
    on organization_plans (organization_id);


create table if not exists plan_features
(
    id           serial
        primary key,
    plan_id      int          not null,
    feature_name varchar(255) not null,
    constraint plan_features_ibfk_1
        foreign key (plan_id) references plans (id)
);

create index plan_id_idx
    on plan_features (plan_id);

create table if not exists rides
(
    id         serial
        primary key,
    branch_id  int                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    ended_at   timestamp                           null,
    constraint ride_branches_id_fk
        foreign key (branch_id) references branches (id)
);

create table if not exists rides_to_orders
(
    id       serial
        primary key,
    ride_id  int not null,
    order_id int not null,
    constraint ride_to_orders_ride_id_fk
        foreign key (ride_id) references rides (id),
    constraint rides_to_orders_orders_id_fk
        foreign key (order_id) references orders (id)
);

create table if not exists roles
(
    id          serial
        primary key,
    name        varchar(64)  not null,
    description varchar(128) null
);

CREATE TABLE IF NOT EXISTS users (
                                     id serial
                                         PRIMARY KEY,
                                     login VARCHAR(50) NOT NULL UNIQUE,
                                     password TEXT NOT NULL,
                                     name VARCHAR(50),
                                     surname VARCHAR(50),
                                     patronymic VARCHAR(50),
                                     organization_id INT NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                     CONSTRAINT users_organizations_id_fk FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

create table if not exists user_roles
(
    id      serial
        primary key,
    role_id int not null,
    user_id int not null,
    constraint user_roles_role_id_fk
        foreign key (role_id) references roles (id),
    constraint user_roles_users_id_fk
        foreign key (user_id) references users (id)
);

-- Trigger function to update `updated_at` timestamp on row update
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function on update
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Создание таблицы сессий
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    refresh_token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT sessions_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Индексы для ускорения поиска по user_id и session_token
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_session_token ON sessions (session_token);
CREATE INDEX idx_sessions_refresh_token ON sessions (refresh_token);

-- Триггер для автоматического удаления просроченных сессий
CREATE OR REPLACE FUNCTION delete_expired_sessions()
    RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_expired_sessions
    BEFORE INSERT ON sessions
    FOR EACH ROW
EXECUTE FUNCTION delete_expired_sessions();