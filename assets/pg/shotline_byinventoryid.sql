SELECT sl.shotline_id,
  sl.name,
  sl.alt_names,
  sl.year,
  sl.remark,
  sp.shotpoint_id,
  sp.shotpoint_number
FROM shotline AS sl
JOIN shotpoint AS sp
  ON sp.shotline_id = sl.shotline_id
JOIN inventory_shotpoint AS isp
  ON isp.shotpoint_id = sp.shotpoint_id
WHERE isp.inventory_id = $1
ORDER BY sl.shotline_id, sp.shotpoint_number
