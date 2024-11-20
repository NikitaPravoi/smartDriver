create table orders
(
    id            bigint generated always as identity
        primary key,
    customer_name text      not null,
    city          text,
    street        text,
    apartment     text,
    floor         integer,
    entrance      integer,
    comment       text,
    cost          numeric   not null,
    status        text,
    location      point     not null,
    created_at    timestamp not null
);

create table organizations
(
    id             bigint generated always as identity
        primary key,
    name           text              not null,
    balance        numeric default 0 not null,
    iiko_api_token text              not null
        constraint organizations_iiko_id_key
            unique
);

create table branches
(
    id              bigint generated always as identity
        primary key,
    name            text   not null,
    location        point  not null,
    organization_id bigint not null
        references organizations
);

create table plans
(
    id             bigint generated always as identity
        primary key,
    name           text           not null,
    cost           numeric(10, 2) not null,
    employee_limit integer        not null
);

create table organization_plans
(
    id              bigint generated always as identity
        primary key,
    organization_id bigint not null
        references organizations,
    plan_id         bigint not null
        references plans,
    start_date      date   not null,
    end_date        date
);

create index organization_id_idx
    on organization_plans (organization_id);

create table plan_features
(
    id           bigint generated always as identity
        primary key,
    plan_id      bigint not null
        references plans,
    feature_name text   not null
);

create index plan_id_idx
    on plan_features (plan_id);

create table rides
(
    id         bigint generated always as identity
        primary key,
    branch_id  bigint                              not null
        references branches,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    ended_at   timestamp
);

create table rides_to_orders
(
    id       bigint generated always as identity
        primary key,
    ride_id  bigint not null
        references rides,
    order_id bigint not null
        references orders
);

create table roles
(
    id          bigint generated always as identity
        primary key,
    name        text not null,
    description text
);

create table users
(
    id              bigint generated always as identity
        primary key,
    login           text                                not null
        unique,
    password        text                                not null,
    name            text,
    surname         text,
    patronymic      text,
    organization_id bigint                              not null
        references organizations,
    created_at      timestamp default CURRENT_TIMESTAMP not null,
    updated_at      timestamp default CURRENT_TIMESTAMP not null
);

create table user_roles
(
    id      bigint generated always as identity
        primary key,
    role_id bigint not null
        references roles,
    user_id bigint not null
        references users
);

create table sessions
(
    id            bigint generated always as identity
        primary key,
    user_id       bigint                              not null
        references users
            on delete cascade,
    session_token text                                not null
        unique,
    refresh_token text                                not null
        unique,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    expires_at    timestamp                           not null
);

create index idx_sessions_user_id
    on sessions (user_id);

create index idx_sessions_session_token
    on sessions (session_token);

create index idx_sessions_refresh_token
    on sessions (refresh_token);

create table payments
(
    id              bigint generated always as identity
        primary key,
    organization_id bigint                              not null
        references organizations,
    amount          numeric(10, 2)                      not null,
    payment_date    timestamp default CURRENT_TIMESTAMP not null,
    payment_method  text                                not null
);

create table revision
(
    revision_id     bigint not null
        constraint revision_id
            primary key,
    organization_id bigint
        constraint revision_organizations_id_fk
            references organizations
);

create unique index revision_organization_id_uindex
    on revision (organization_id);

