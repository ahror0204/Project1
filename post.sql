CREATE TABLE IF NOT EXISTS posts
(
    id uuid not null primary key,
    user_id uuid,
    name text,
    media text[],
    description text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS medias(
    id uuid primary key not null,
    type text ,
    link text,
    post_id uuid,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);