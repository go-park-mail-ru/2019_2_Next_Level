drop table if exists Session;
drop table if exists Receiver;
drop table if exists Folder;
drop table if exists Message;
drop table if exists Users;

drop type if exists sex_type;
drop type if exists mail_direction;
create type sex_type as ENUM ('male', 'female');
create type mail_direction as ENUM ('in', 'out');

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
    birthdate date NOT NULL DEFAULT NOW()
);

create table if not exists Session
(
    login CITEXT NOT NULL REFERENCES Users (login),
    token uuid NOT NULL
);

create table if not exists Folder (
    name TEXT NOT NULL DEFAULT 'inbox',
    owner citext REFERENCES Users (login) NOT NULL,
    count integer DEFAULT 0,
    UNIQUE (name, owner)
);

create table if not exists Message
(
    id bigserial not null primary key,

    sender CITEXT not null,
    subject TEXT NOT NULL DEFAULT '',
    body text default '' not null,
    direction mail_direction NOT NULL DEFAULT 'in',

    time timestamp with time zone not null default NOW(),
    folder TEXT DEFAULT 'inbox' NOT NULL,
    isRead bool DEFAULT false,
    isMarked bool DEFAULT false
);


create table if not exists Receiver
(
    id bigserial not null primary key,
    mailid bigint references Message (id) on delete cascade,
    email citext not null
);

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

-- Обновляет счетчики у папок
CREATE OR REPLACE FUNCTION folder_counter() RETURNS trigger AS $folder_counter$
    BEGIN
        if tg_op='INSERT' OR tg_op='UPDATE' then
            if NEW.direction='in' then
                UPDATE Folder SET count=count+1 WHERE owner IN (SELECT Receiver.email FROM Receiver WHERE mailid=NEW.id) AND name=NEW.folder;
                RAISE NOTICE 'IN';
            else
                UPDATE Folder SET count=count+1 WHERE owner=NEW.sender AND name=NEW.folder;
                RAISE NOTICE 'OUT';
            end if;
        end if;
        if tg_op='DELETE' OR tg_op='UPDATE' then
            if OLD.direction='in' then
                UPDATE Folder SET count=count-1 WHERE owner IN (SELECT Receiver.email FROM Receiver WHERE mailid=OLD.id) AND name=OLD.folder;
            else
                UPDATE Folder SET count=count-1 WHERE owner=OLD.sender AND name=OLD.folder;
            end if;
            RETURN OLD;
        end if;
        RETURN NEW;
    END
$folder_counter$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS folder_counter ON Message;
CREATE CONSTRAINT TRIGGER folder_counter AFTER INSERT OR DELETE OR UPDATE ON Message
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE PROCEDURE folder_counter();

INSERT INTO Users (login, password, sault, firstname, secondname)
    VALUES ('admin', 'wedewde', 'wedewdewd', 'Ian', 'Ivanov');
BEGIN;
INSERT INTO Message (sender, subject, body, direction, folder)
    VALUES ('aa@mail.ru', 'Test message', 'Test Body', 'in', 'inbox'),
            ('admin', 'Test outcoming', 'Test body', 'out', 'sent');
INSERT INTO Receiver (mailid, email) VALUES (1, 'admin'), (2, 'aa@mail.ru');
COMMIT;
