-- name: CreatePosts :one
insert into posts(
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
) returning *;

-- name: GetPosts :many
select posts.title, posts.url, posts.description,
posts.published_at, feeds.name as feed_name
from posts
inner join feeds on feeds.id = posts.feed_id
inner join feeds_follows on posts.feed_id = feeds_follows.feed_id
where feeds_follows.user_id = $1
order by posts.published_at desc
limit $2;
