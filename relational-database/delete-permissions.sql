CREATE OR REPLACE FUNCTION delete_folder_permission(
    _account_id VARCHAR,
    _folder_id VARCHAR
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
    _primary_folder_files record;
    _child_folder_files record;
    _child_folders record;
BEGIN
    -- Delete permission from selected folder
    DELETE FROM account_folder
    WHERE account_id = _account_id::UUID
    AND folder_id = _folder_id::UUID;

    -- Delete permissions from files in folder
    FOR _primary_folder_files IN 
    (
        SELECT id FROM file WHERE folder_id = _folder_id::UUID
    )
    LOOP
        DELETE FROM account_file
        WHERE account_id = _account_id::UUID
        AND file_id = _primary_folder_files.id;
    END LOOP;

    -- Loop to remove permissions from all child folders
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
        -- Delete permission from child folder
        DELETE FROM account_folder
        WHERE account_id = _account_id::UUID
        AND folder_id = _child_folders._child_id;

        -- Delete permissions from files in child folder
        FOR _child_folder_files IN 
        (
            SELECT id FROM file WHERE folder_id = _child_folders._child_id
        )
        LOOP
            DELETE FROM account_file
            WHERE account_id = _account_id::UUID
            AND file_id = _child_folder_files.id;
        END LOOP;
    END LOOP;

    RETURN true;
END;
$$;

CREATE OR REPLACE FUNCTION delete_account_file_permission(
    _account_id VARCHAR,
    _file_id VARCHAR,
    permission SMALLINT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM account_file
    WHERE account_id = _account_id::UUID
    AND file_id = _file_id::UUID;

    RETURN true;
END;
$$;