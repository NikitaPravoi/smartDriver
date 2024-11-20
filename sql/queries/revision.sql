-- name: CreateRevision :exec
INSERT INTO revision (revision_id, organization_id) VALUES ($1, $2);

-- name: GetLastRevision :one
SELECT
    revision_id
FROM revision
WHERE organization_id = $1
ORDER BY revision_id DESC
LIMIT 1;