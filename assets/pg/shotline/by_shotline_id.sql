SELECT sl.shotline_id,
    sl.name,
    sl.alt_names,
    sl.year,
    sl.remark,
        MIN(sp.shotpoint_number) as shotpoint_min,
        MAX(sp.shotpoint_number) as shotpoint_max
FROM shotline AS sl
JOIN shotpoint AS sp
    ON sp.shotline_id = sl.shotline_id
WHERE sl.shotline_id = $1
GROUP BY sl.shotline_id
