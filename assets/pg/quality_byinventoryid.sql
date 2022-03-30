SELECT inventory_quality_id,
  check_date,
  remark,
  username,
  ARRAY_TO_JSON(issues) AS issues
FROM inventory_quality
WHERE inventory_id = $1
ORDER BY check_date DESC
