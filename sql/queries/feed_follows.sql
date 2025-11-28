-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feeds_follows(id, created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4, $5)
    returning *
)
select
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
    from inserted_feed_follow
inner join users on user_id = users.id
inner join feeds on feed_id = feeds.id;


-- name: GetFeedFollowsForUser :many

select users.name as user_name, feeds.name as feed_name
from
feeds_follows
inner join users on user_id = users.id
inner join feeds on feed_id = feeds.id
where users.name = $1;
