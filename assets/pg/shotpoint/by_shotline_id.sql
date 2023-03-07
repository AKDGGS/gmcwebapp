SELECT sp.shotline_id AS "ID", sp.shotpoint_number AS "Number"
FROM shotpoint AS sp
WHERE sp.shotline_id = $1
ORDER BY shotpoint_number;
