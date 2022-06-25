SELECT path_cache AS "Container", COUNT(*) AS "Count"
FROM container
WHERE active
GROUP BY path_cache
HAVING COUNT(*) > 1
