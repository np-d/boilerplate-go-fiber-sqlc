CREATE
    or replace function update_updated_at() returns trigger
    language plpgsql AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE TABLE if not exists users
(
    id
               serial
        primary
            key,
    display_name
               varchar(50)                         not null,
    username   varchar(30)                         not null
        constraint users_username_uk unique,
    email      varchar(100)                        not null
        constraint users_email_uk unique,
    password   varchar(255)                        not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);
CREATE
    or replace trigger trigger_update_users_updated_at
    after
        update
    ON users
    for each row
execute procedure update_updated_at();