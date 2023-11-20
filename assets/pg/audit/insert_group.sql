INSERT INTO audit_group (remark) VALUES ($1)
RETURNING audit_group_id;
