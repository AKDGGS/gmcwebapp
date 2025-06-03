-- If the container parent_id changed, update the path_cache
CREATE OR REPLACE FUNCTION container_parent_path_cache_fn()
RETURNS TRIGGER AS $$
BEGIN
	NEW.path_cache = COALESCE((
			SELECT path_cache FROM container
			WHERE container_id = NEW.parent_container_id
	) || '/', '') || NEW.name;
	RETURN NEW;
END; $$ LANGUAGE 'plpgsql';

DROP TRIGGER IF EXISTS container_parent_path_cache_tr ON container;
CREATE TRIGGER container_parent_path_cache_tr
BEFORE INSERT OR UPDATE ON container
FOR EACH ROW EXECUTE PROCEDURE container_parent_path_cache_fn();

-- If the container parent_container_id changes, or the name changes
-- update all the container's children with the new path
CREATE OR REPLACE FUNCTION container_children_path_cache_fn()
RETURNS TRIGGER AS $$
BEGIN
	IF (NEW.parent_container_id <> OLD.parent_container_id) OR (NEW.name <> OLD.name) THEN
		WITH RECURSIVE t AS ((
			SELECT 0 AS depth, c.container_id, c.name::text
			FROM container AS c
			WHERE c.parent_container_id = NEW.container_id
		) UNION ALL (
			SELECT depth +1 AS depth, c.container_id,
				(t.name || '/' || c.name)::text AS name
			FROM container AS c
			JOIN t ON c.parent_container_id = t.container_id
			WHERE depth <= 20
		))
		UPDATE container SET
			path_cache = NEW.path_cache || '/' || t.name
		FROM t
		WHERE t.container_id = container.container_id;
	END IF;
	RETURN NEW;
END; $$ LANGUAGE 'plpgsql';

DROP TRIGGER IF EXISTS container_children_path_cache_tr ON container;
CREATE TRIGGER container_children_path_cache_tr AFTER UPDATE ON container
FOR EACH ROW EXECUTE PROCEDURE container_children_path_cache_fn();

-- Create function/trigger for inventory change logging
CREATE OR REPLACE FUNCTION inventory_container_log_fn()
RETURNS TRIGGER AS $$
BEGIN
	IF TG_OP = 'INSERT' THEN
		IF NEW.container_id IS NOT NULL THEN
			INSERT INTO inventory_container_log (
				inventory_id, destination
			) VALUES (
				NEW.inventory_id, (
					SELECT path_cache FROM container
					WHERE container_id = NEW.container_id
				)
			);
		END IF;
	ELSIF TG_OP = 'UPDATE' THEN
		IF COALESCE(OLD.container_id, 0) <> COALESCE(NEW.container_id, 0) THEN
			INSERT INTO inventory_container_log (
				inventory_id, destination
			) VALUES (
				NEW.inventory_id, (
					SELECT path_cache FROM container
					WHERE container_id = NEW.container_id
				)
			);
		END IF;
	END IF;

	RETURN NEW;
END; $$ LANGUAGE 'plpgsql';

DROP TRIGGER IF EXISTS inventory_container_log_tr ON inventory;
CREATE TRIGGER inventory_container_log_tr
AFTER INSERT OR UPDATE ON inventory
FOR EACH ROW EXECUTE PROCEDURE inventory_container_log_fn();


-- Create function/trigger for container change logging
CREATE OR REPLACE FUNCTION container_log_fn()
RETURNS TRIGGER AS $$
BEGIN
	IF TG_OP = 'INSERT' THEN
		IF NEW.parent_container_id IS NOT NULL THEN
			INSERT INTO container_log (
				container_id, destination
			) VALUES (
				NEW.container_id, (
					SELECT path_cache FROM container
					WHERE container_id = NEW.parent_container_id
				)
			);
		END IF;
	ELSIF TG_OP = 'UPDATE' THEN
		IF COALESCE(OLD.parent_container_id, 0) <> COALESCE(NEW.parent_container_id, 0) THEN
			INSERT INTO container_log (
				container_id, destination
			) VALUES (
				NEW.container_id, (
					SELECT path_cache FROM container
					WHERE container_id = NEW.parent_container_id
				)
			);
		END IF;
	END IF;

	RETURN NEW;
END; $$ LANGUAGE 'plpgsql';

DROP TRIGGER IF EXISTS container_log_tr ON container;
CREATE TRIGGER container_log_tr
AFTER INSERT OR UPDATE ON container
FOR EACH ROW EXECUTE PROCEDURE container_log_fn();


-- Create function for modified date touching on update/insert
CREATE OR REPLACE FUNCTION modified_date_fn()
RETURNS TRIGGER AS $$
BEGIN
	NEW.modified_date = NOW();
	RETURN NEW;
END; $$ language 'plpgsql';

-- Create function for modified user touching on update/insert
CREATE OR REPLACE FUNCTION modified_user_fn()
RETURNS TRIGGER AS $$
BEGIN
	IF session_user <> 'gmc_app' OR NEW.modified_user IS NULL THEN
		NEW.modified_user = session_user;
	END IF;
	RETURN NEW;
END; $$ language 'plpgsql';

-- Set trigger for inventory modified date
DROP TRIGGER IF EXISTS inventory_modified_date_tr ON inventory;
CREATE TRIGGER inventory_modified_date_tr BEFORE INSERT OR UPDATE ON inventory
FOR EACH ROW EXECUTE PROCEDURE modified_date_fn();

-- Set trigger for inventory modified user
DROP TRIGGER IF EXISTS inventory_modified_user_tr ON inventory;
CREATE TRIGGER inventory_modified_user_tr BEFORE INSERT OR UPDATE ON inventory
FOR EACH ROW EXECUTE PROCEDURE modified_user_fn();

-- Set trigger for outcrop modified date
DROP TRIGGER IF EXISTS outcrop_modified_date_tr ON outcrop;
CREATE TRIGGER outcrop_modified_date_tr BEFORE INSERT OR UPDATE ON outcrop
FOR EACH ROW EXECUTE PROCEDURE modified_date_fn();

-- Set trigger for outcrop modified user
DROP TRIGGER IF EXISTS outcrop_modified_user_tr ON outcrop;
CREATE TRIGGER outcrop_modified_user_tr BEFORE INSERT OR UPDATE ON outcrop
FOR EACH ROW EXECUTE PROCEDURE modified_user_fn();

-- Set trigger for borehole modified date
DROP TRIGGER IF EXISTS borehole_modified_date_tr ON borehole;
CREATE TRIGGER borehole_modified_date_tr BEFORE INSERT OR UPDATE ON borehole
FOR EACH ROW EXECUTE PROCEDURE modified_date_fn();

-- Set trigger for borehole modified user
DROP TRIGGER IF EXISTS borehole_modified_user_tr ON borehole;
CREATE TRIGGER borehole_modified_user_tr BEFORE INSERT OR UPDATE ON borehole
FOR EACH ROW EXECUTE PROCEDURE modified_user_fn();

-- Set trigger for well modified date
DROP TRIGGER IF EXISTS well_modified_date_tr ON well;
CREATE TRIGGER well_modified_date_tr BEFORE INSERT OR UPDATE ON well
FOR EACH ROW EXECUTE PROCEDURE modified_date_fn();

-- Set trigger for well modified user
DROP TRIGGER IF EXISTS well_modified_user_tr ON well;
CREATE TRIGGER well_modified_user_tr BEFORE INSERT OR UPDATE ON well
FOR EACH ROW EXECUTE PROCEDURE modified_user_fn();
