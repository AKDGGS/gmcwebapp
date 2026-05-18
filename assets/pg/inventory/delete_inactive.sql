WITH inactive AS (
	SELECT inventory_id FROM inventory WHERE active = false
),

delete_borehole AS (
    DELETE FROM inventory_borehole 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_collector AS (
    DELETE FROM inventory_collector 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_container_log AS (
    DELETE FROM inventory_container_log 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_file AS (
    DELETE FROM inventory_file 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_publication AS (
    DELETE FROM inventory_publication
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_note AS (
    DELETE FROM inventory_note 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_outcrop AS (
    DELETE FROM inventory_outcrop 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_quality AS (
    DELETE FROM inventory_quality 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_shotepoint AS (
    DELETE FROM inventory_shotpoint 
    WHERE inventory_id 
    IN (SELECT inventory_id inactive)
),

delete_url AS (
    DELETE FROM inventory_url 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
),

delete_well AS (
    DELETE FROM inventory_well 
    WHERE inventory_id 
    IN (SELECT inventory_id FROM inactive)
)

DELETE FROM inventory WHERE active=false;