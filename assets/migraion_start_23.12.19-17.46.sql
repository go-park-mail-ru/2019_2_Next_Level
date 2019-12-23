ALTER TABLE Session ADD COLUMN expires date NOT NULL DEFAULT NOW() + INTERVAL '1 DAY';

CREATE OR REPLACE FUNCTION receiver_own() RETURNS trigger AS $receiver_own$
BEGIN
    if NEW.email LIKE '%nl-mail.ru' then
        NEW.email = (SELECT substr(NEW.email, 0, position('@' in NEW.email)));
    end if;
    return NEW;
END
$receiver_own$ LANGUAGE plpgsql;