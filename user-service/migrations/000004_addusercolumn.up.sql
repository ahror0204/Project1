ALTER TABLE users ADD COLUMN password VARCHAR(100) UNIQUE NOT NULL;
ALTER TABLE users ADD COLUMN email_code TEXT;

