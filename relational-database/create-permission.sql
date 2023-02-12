CREATE OR REPLACE FUNCTION create_folder_permission(
    _account_id VARCHAR,
    _folder_id VARCHAR,
    _permission SMALLINT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
    _primary_folder_files record;
    _child_folder_files record;
    _child_folders record;
BEGIN
    -- Add permission to selected folder
    INSERT INTO account_folder(account_id, folder_id, permission) 
    VALUES(_account_id::UUID, _folder_id::UUID, _permission);

    -- Add permissions to files in folder
    FOR _primary_folder_files IN 
    (
        SELECT id FROM file WHERE folder_id = _folder_id::UUID
    )
    LOOP
        INSERT INTO account_file(account_id, file_id, permission)
        VALUES(_account_id::UUID, _primary_folder_files.id, _permission);
    END LOOP;

    -- Loop to add permissions to all child folders
    FOR _child_folders IN
    (
        -- Find all child folders
        WITH RECURSIVE cte(_child_id) as 
        (
            SELECT id
            FROM folder
            WHERE parent_id = _folder_id::UUID

            UNION ALL

            SELECT folder.id
            FROM folder, cte
            WHERE folder.parent_id = cte._child_id
        )
        SELECT _child_id FROM cte
    )
    LOOP
        -- Add permission to child folder
        INSERT INTO account_folder(account_id, folder_id, permission)
        VALUES(_account_id::UUID, _child_folders._child_id, _permission);

        -- Add permissions to files in child folder
        FOR _child_folder_files IN 
        (
            SELECT id FROM file WHERE folder_id = _child_folders._child_id
        )
        LOOP
            INSERT INTO account_file(account_id, file_id, permission)
            VALUES(_account_id::UUID, _child_folder_files.id, _permission);
        END LOOP;
    END LOOP;

    RETURN true;
END;
$$;

CREATE OR REPLACE FUNCTION create_file_permission(
    _account_id VARCHAR,
    _file_id VARCHAR,
    _permission SMALLINT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO account_file(account_id, file_id, permission)
    VALUES(_account_id::UUID, _file_id::UUID, _permission);

    RETURN true;
END;
$$;

CREATE OR REPLACE FUNCTION edit_folder_permission(
    _account_id VARCHAR,
    _folder_id VARCHAR,
    _permission SMALLINT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
    _primary_folder_files record;
    _child_folder_files record;
    _child_folders record;
BEGIN
    -- Add permission to selected folder
    UPDATE account_folder
    SET permission = _permission
    WHERE account_id = _account_id::UUID
    AND folder_id = _folder_id::UUID;

    -- Add permissions to files in folder
    FOR _primary_folder_files IN 
    (
        SELECT id FROM file WHERE folder_id = _folder_id::UUID
    )
    LOOP
        UPDATE account_file
        SET permission = _permission
        WHERE account_id = _account_id::UUID
        AND file_id = _primary_folder_files.id;
    END LOOP;

    -- Loop to add permissions to all child folders
    FOR _child_folders IN
    (
        -- Find all child folders
        WITH RECURSIVE cte(_child_id) as 
        (
            SELECT id
            FROM folder
            WHERE parent_id = _folder_id::UUID

            UNION ALL

            SELECT folder.id
            FROM folder, cte
            WHERE folder.parent_id = cte._child_id
        )
        SELECT _child_id FROM cte
    )
    LOOP
        -- Add permission to child folder
        UPDATE account_folder
        SET permission = _permission
        WHERE account_id = _account_id::UUID
        AND folder_id = _child_folders._child_id;

        -- Add permissions to files in child folder
        FOR _child_folder_files IN 
        (
            SELECT id FROM file WHERE folder_id = _child_folders._child_id
        )
        LOOP
            UPDATE account_file
            SET permission = _permission
            WHERE account_id = _account_id::UUID
            AND file_id = _child_folder_files.id;
        END LOOP;
    END LOOP;

    RETURN true;
END;
$$;

CREATE OR REPLACE FUNCTION edit_file_permission(
    _account_id VARCHAR,
    _file_id VARCHAR,
    _permission SMALLINT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE account_file
    SET permission = _permission
    WHERE account_id = _account_id::UUID
    AND file_id = _file_id::UUID;

    RETURN true;
END;
$$;