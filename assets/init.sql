drop table if exists Session;
drop table if exists Receiver;
drop table if exists Message;
drop table if exists Folder;
drop table if exists ChatMessage;
drop table if exists Chat;
drop table if exists Users;

drop type if exists sex_type;
drop type if exists mail_direction;
drop type if exists user_role;
create type sex_type as ENUM ('male', 'female');
create type mail_direction as ENUM ('in', 'out');
create type user_role as ENUM ('common', 'support', 'admin');

CREATE EXTENSION IF NOT EXISTS CITEXT WITH SCHEMA public;

create table if not exists Users
(
    login CITEXT PRIMARY KEY NOT NULL,
    password bytea not null,
    sault bytea not null,

    firstname TEXT NOT NULL DEFAULT '',
    secondname TEXT NOT NULL DEFAULT '',
    sex sex_type DEFAULT 'male' NOT NULL,
    avatar TEXT default 'default.png' NOT NULL,
    birthdate date NOT NULL DEFAULT NOW(),
    role user_role NOT NULL DEFAULT 'common'
);

create table if not exists Session
(
    login CITEXT NOT NULL REFERENCES Users (login) ON DELETE CASCADE,
    token uuid NOT NULL
);

create table if not exists Folder (
    id BIGSERIAL not null,
    name TEXT NOT NULL DEFAULT 'inbox',
    owner citext REFERENCES Users (login) ON DELETE CASCADE NOT NULL ,
    count integer DEFAULT 0,
    UNIQUE (name, owner)
);

create table if not exists Message
(
    id bigserial not null primary key,
    owner citext REFERENCES Users (login) ON DELETE CASCADE,

    sender CITEXT not null,
    subject TEXT NOT NULL DEFAULT '',
    body text default '' not null,
    direction mail_direction NOT NULL DEFAULT 'in',

    time timestamp with time zone not null default NOW(),
    folder TEXT DEFAULT 'inbox' NOT NULL,
    isRead bool DEFAULT false,
    isMarked bool DEFAULT false,
    CONSTRAINT cnst FOREIGN KEY (folder, owner) REFERENCES Folder(name, owner)
);


create table if not exists Receiver
(
    id bigserial not null primary key,
    mailid bigint references Message (id) on delete cascade,
    email citext not null
);

create table if not exists Chat(
    id bigserial not null primary key,
    userNick citext REFERENCES Users(login) not null,
    supportNick citext REFERENCES Users(login) not null,
    startDate timestamptz NOT NULL DEFAULT now(),
    isOpen bool NOT NULL DEFAULT true,
    theme text NOT NULL DEFAULT ''
);

create table if not exists ChatMessage(
    id bigserial not null primary key,
    chatId bigserial References Chat (id) NOT NULL,
    sent timestamptz not null default now(),
    wasRead bool NOT NULL DEFAULT false,
    body text not null DEFAULT '',
    author citext REFERENCES Users(login) NOT NULL

);

CREATE OR REPLACE FUNCTION get_chat_field(userid_ citext) RETURNS text as $get_chat_id$
    DECLARE
        role_ user_role;
    BEGIN
        role_ := (SELECT role FROM Users WHERE login=userid_);
        IF role_='common' then
            RETURN 'userNick';
        ELSE
            RETURN 'supportNick';
        end if;

    end;
$get_chat_id$ LANGUAGE plpgsql;


-- Триггеры

-- Создаёт каждому пользователю стандартные папки
CREATE OR REPLACE FUNCTION user_created() RETURNS trigger AS $user_created$
    BEGIN
        INSERT INTO Folder (name, owner) VALUES
            ('inbox', NEW.login),
            ('sent', NEW.login),
            ('proceed', NEW.login),
            ('spam', NEW.login),
            ('trash', NEW.login);
        RETURN NEW;
    END
$user_created$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS user_created ON Users;
CREATE TRIGGER user_created AFTER INSERT ON Users
    FOR EACH ROW EXECUTE PROCEDURE user_created();


-- CREATE OR REPLACE FUNCTION user_changed() RETURNS trigger AS $user_changed$
--     BEGIN
--         IF NEW.avatars=='' then
--             NEW.avatars = 'default';
--         end if;
--         return NEW;
--     END
-- $user_changed$ LANGUAGE plpgsql;
-- DROP TRIGGER IF EXISTS user_changed ON Users;
-- CREATE TRIGGER user_changed AFTER UPDATE OR INSERT ON Users
--     FOR EACH ROW EXECUTE PROCEDURE user_changed();


-- Обновляет счетчики у папок
CREATE OR REPLACE FUNCTION folder_counter() RETURNS trigger AS $folder_counter$
    BEGIN
        if tg_op='INSERT' OR tg_op='UPDATE' then
            UPDATE Folder SET count=count+1 WHERE owner=NEW.owner AND name=NEW.folder;
        end if;
        if tg_op='DELETE' OR tg_op='UPDATE' then
            UPDATE Folder SET count=count-1 WHERE owner=OLD.owner AND name=OLD.folder;
            RETURN OLD;
        end if;
        RETURN NEW;
    END
$folder_counter$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS folder_counter ON Message;
CREATE CONSTRAINT TRIGGER folder_counter AFTER INSERT OR DELETE OR UPDATE ON Message
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE PROCEDURE folder_counter();

-- Пзаполнение поля "owner" для совместимости со старыми запросами
CREATE OR REPLACE FUNCTION message_owner() RETURNS trigger AS $message_owner$
    BEGIN
        if NEW.owner IS NULL then
            if NEW.sender IN (SELECT login FROM Users) then
                NEW.owner = NEW.sender;
            else
                NEW.owner = (SELECT email FROM receiver WHERE mailid=NEW.id LIMIT 1);
            end if;
            RETURN NEW;
        end if;
        RETURN NEW;
    END
$message_owner$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS message_owner ON Message;
CREATE TRIGGER message_owner BEFORE INSERT ON Message
    FOR EACH ROW EXECUTE PROCEDURE message_owner();

-- Перемещает сообщения из папок при удалении
CREATE or replace function on_remove_folder() returns trigger as $on_remove_folder$
    begin
--         CREATE VIEW list AS SELECT id from Message where folder=OLD.name;
--         UPDATE Message SET folder='inbox' WHERE id in (SELECT id from list) AND direction='in';
--         UPDATE Message SET folder='sent'  WHERE id in (SELECT id from list) AND direction='out';
        UPDATE Message SET folder=CASE(direction)
                                    WHEN 'in' THEN 'inbox'
                                    WHEN 'out'THEN 'sent'
                                    END
                    WHERE id in (SELECT id from Message where folder=OLD.name) AND direction='out';
        return old;
    end
$on_remove_folder$ language plpgsql;
drop trigger if exists on_remove_folder ON Folder;
create trigger on_remove_folder BEFORE DELETE ON Folder
    FOR EACH row execute procedure on_remove_folder();

INSERT INTO Users (login, password, sault, firstname, secondname)
    VALUES ('admin', 'wedewde', 'wedewdewd', 'Ian', 'Ivanov');
BEGIN;
INSERT INTO Message (sender, subject, body, direction, folder)
    VALUES ('aa@mail.ru', 'Test message', 'Test Body', 'in', 'inbox'),
            ('admin', 'Test outcoming', 'Test body', 'out', 'sent');
INSERT INTO Receiver (mailid, email) VALUES (1, 'admin'), (2, 'aa@mail.ru');
COMMIT;
