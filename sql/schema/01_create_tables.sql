create table if not exists orders (
                                      id bigint primary key generated always as identity,
                                      customer_name text not null,
                                      customer_phone text,
                                      city text,
                                      street text,
                                      apartment text,
                                      floor int,
                                      entrance int,
                                      comment text,
                                      cost numeric(10, 2) not null,
                                      status int default 0,
                                      created_at timestamp not null
);

create table if not exists organizations (
                                             id bigint primary key generated always as identity,
                                             name text not null,
                                             balance real default 0 not null,
                                             iiko_id text not null unique
);

create table if not exists branches (
                                        id bigint primary key generated always as identity,
                                        name text not null,
                                        location point not null,
                                        organization_id bigint not null,
                                        foreign key (organization_id) references organizations (id)
);

create table if not exists plans (
                                     id bigint primary key generated always as identity,
                                     name text not null,
                                     cost numeric(10, 2) not null,
                                     employee_limit int not null
);

create table if not exists organization_plans (
                                                  id bigint primary key generated always as identity,
                                                  organization_id bigint not null,
                                                  plan_id bigint not null,
                                                  start_date date not null,
                                                  end_date date,
                                                  foreign key (organization_id) references organizations (id),
                                                  foreign key (plan_id) references plans (id)
);

create index organization_id_idx on organization_plans using btree (organization_id);

create table if not exists plan_features (
                                             id bigint primary key generated always as identity,
                                             plan_id bigint not null,
                                             feature_name text not null,
                                             foreign key (plan_id) references plans (id)
);

create index plan_id_idx on plan_features using btree (plan_id);

create table if not exists rides (
                                     id bigint primary key generated always as identity,
                                     branch_id bigint not null,
                                     created_at timestamp default current_timestamp not null,
                                     ended_at timestamp,
                                     foreign key (branch_id) references branches (id)
);

create table if not exists rides_to_orders (
                                               id bigint primary key generated always as identity,
                                               ride_id bigint not null,
                                               order_id bigint not null,
                                               foreign key (ride_id) references rides (id),
                                               foreign key (order_id) references orders (id)
);

create table if not exists roles (
                                     id bigint primary key generated always as identity,
                                     name text not null,
                                     description text
);

create table if not exists users (
                                     id bigint primary key generated always as identity,
                                     login text not null unique,
                                     password text not null,
                                     name text,
                                     surname text,
                                     patronymic text,
                                     organization_id bigint not null,
                                     created_at timestamp default current_timestamp not null,
                                     updated_at timestamp default current_timestamp not null,
                                     foreign key (organization_id) references organizations (id)
);

create table if not exists user_roles (
                                          id bigint primary key generated always as identity,
                                          role_id bigint not null,
                                          user_id bigint not null,
                                          foreign key (role_id) references roles (id),
                                          foreign key (user_id) references users (id)
);

create
    or replace function update_updated_at_column () returns trigger as $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

create trigger update_users_updated_at before
    update on users for each row
execute function update_updated_at_column ();

create table if not exists sessions (
                                        id bigint primary key generated always as identity,
                                        user_id bigint not null,
                                        session_token text not null unique,
                                        refresh_token text not null unique,
                                        created_at timestamp default current_timestamp not null,
                                        expires_at timestamp not null,
                                        foreign key (user_id) references users (id) on delete cascade
);

create index idx_sessions_user_id on sessions using btree (user_id);

create index idx_sessions_session_token on sessions using btree (session_token);

create index idx_sessions_refresh_token on sessions using btree (refresh_token);

create
    or replace function delete_expired_sessions () returns trigger as $$
BEGIN
    DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

create trigger trigger_delete_expired_sessions before insert on sessions for each row
execute function delete_expired_sessions ();

create table if not exists payments (
                                        id bigint primary key generated always as identity,
                                        organization_id bigint not null,
                                        amount numeric(10, 2) not null,
                                        payment_date timestamp default current_timestamp not null,
                                        payment_method text not null,
                                        foreign key (organization_id) references organizations (id)
);

create
    or replace function withdraw_balance (o_id bigint, amount numeric) returns void as $$
DECLARE
    current_balance numeric;
BEGIN
    -- Get the current balance of the organization
    SELECT balance + COALESCE(SUM(p.amount), 0) INTO current_balance
    FROM organizations o
             LEFT JOIN payments p ON o.id = p.organization_id
    WHERE o.id = o_id
    GROUP BY o.balance;

    -- Check if the balance is sufficient
    IF current_balance >= amount THEN
        -- Insert a payment record to withdraw the amount
        INSERT INTO payments (organization_id, amount, payment_method)
        VALUES (organization_id, -amount, 'withdrawal');
    ELSE
        RAISE NOTICE 'Insufficient balance to withdraw %', amount;
    END IF;
END;
$$ language plpgsql;