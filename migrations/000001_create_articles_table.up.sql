create table if not exists articles (
    id integer not null primary key,
    slug text not null,
    title text not null,
    content text,
    created_at text default (datetime('now'))
);