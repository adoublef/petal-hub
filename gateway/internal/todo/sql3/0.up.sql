create table todos (
    -- uuid_v7
    id text
    -- enforce requiring needing a title
    , title text not null
    , body text 
    , is_complete int
    -- julian
    , updated_at real not null --default julian('now')
    , check (length(title) > 0 and length(title) <= 30)
    , check (length(body) <= 280)
    , check (is_complete in (0,1))
    , primary key (id)
) strict;