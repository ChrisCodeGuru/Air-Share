CREATE OR REPLACE FUNCTION get_folder_parent_permission(
    _account_id VARCHAR,
    _folder_id VARCHAR
)
RETURNS SMALLINT
LANGUAGE plpgsql
AS $$
BEGIN
    SELECT permission
    FROM account_folder
    WHERE account_id = _account_id::UUID
    AND folder_id = (
        SELECT parent_id
        FROM folder
        WHERE id = _folder_id::UUID
    );
END;
$$;

CREATE OR REPLACE FUNCTION get_file_parent_permission(
    _account_id VARCHAR,
    _file_id VARCHAR
)
RETURNS SMALLINT
LANGUAGE plpgsql
AS $$
BEGIN
    SELECT permission
    FROM account_folder
    WHERE account_id = _account_id::UUID
    AND folder_id = (
        SELECT folder_id
        FROM file
        WHERE id = _file_id::UUID
    );
END;
$$;