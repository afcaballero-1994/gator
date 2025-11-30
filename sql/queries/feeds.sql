-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
returning *;

-- name: GetFeeds :many
select feeds.name, feeds.url, users.name as username from feeds join users on user_id = users.id;

-- name: GetFeed :one
select feeds.id from feeds where feeds.url = $1 limit 1;

-- name: MarkFeedAsFetched :exec
update feeds set last_fetched_at = $1, updated_at = $1
where id = $2;

-- name: GetNextFeedToFetch :many
select id, url from feeds order by last_fetched_at asc nulls first;
