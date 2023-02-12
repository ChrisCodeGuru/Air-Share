CREATE OR REPLACE FUNCTION get_owned_contents_folder(
    _account_id VARCHAR
)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT
    --_type UUID
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT folder.id, trim(both '"' from folder.name) "foldername"
    FROM account_folder
    INNER JOIN folder
        ON folder.id = account_folder.folder_id
    WHERE account_folder.account_id::text = _account_id
    AND account_folder.permission = 4
    AND folder.parent_id IS NULL;
END;
$$;

CREATE OR REPLACE FUNCTION get_owned_contents_file(
    _account_id VARCHAR

)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT,
    _sensitive BOOLEAN,
    _hash_value VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT file.id, trim(both '"' from file.name) "filename", file.sensitive, file.hash_value
    FROM account_file
    INNER JOIN file
        ON file.id = account_file.file_id
    WHERE account_file.account_id::text = _account_id
    AND account_file.permission = 4
    AND file.folder_id IS NULL;
END;
$$;

CREATE OR REPLACE FUNCTION get_shared_contents_folder(
    _account_id VARCHAR
)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT,
    _permission SMALLINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT folder.id, trim(both '"' from folder.name) "foldername", if1.permission
    FROM folder
    INNER JOIN account_folder if1
        ON if1.folder_id = folder.id
    WHERE if1.account_id::text = _account_id
    AND if1.permission < 4
    AND (
        SELECT COUNT(*)
        FROM account_folder
        WHERE account_id::text = _account_id
        AND folder_id = folder.parent_id
    ) = 0;
END;
$$;

CREATE OR REPLACE FUNCTION get_shared_contents_file(
    _account_id VARCHAR
)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT,
    _permission SMALLINT,
    _sensitive BOOLEAN,
    _hash_value VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT file.id, trim(both '"' from file.name) "filename", if1.permission, file.sensitive, file.hash_value
    FROM file
    INNER JOIN account_file if1
        ON if1.file_id = file.id
    WHERE if1.account_id::text = _account_id
    AND if1.permission < 4
    AND (
        SELECT COUNT(*)
        FROM account_folder
        WHERE account_id::text = _account_id
        AND folder_id = file.folder_id
    ) = 0;
END;
$$;

CREATE OR REPLACE FUNCTION get_contents_folder(
    _folder_id VARCHAR,
    _account_id VARCHAR
)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT,
    _permission SMALLINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT folder.id, trim(both '"' from name) "foldername", account_folder.permission
    FROM folder
    INNER JOIN account_folder
        ON account_folder.folder_id = folder.id
    WHERE parent_id = _folder_id::UUID
    AND account_folder.account_id = _account_id::UUID;
END;
$$;

CREATE OR REPLACE FUNCTION get_contents_file(
    _folder_id VARCHAR,
    _account_id VARCHAR
)
RETURNS TABLE (
    _item_id UUID,
    _item_name TEXT,
    _permission SMALLINT,
    _sensitive BOOLEAN,
    _hash_value VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT file.id, trim(both '"' from name) "filename", account_file.permission, file.sensitive, file.hash_value
    FROM file
    INNER JOIN account_file
        ON account_file.file_id = file.id
    WHERE folder_id = _folder_id::UUID
    AND account_file.account_id = _account_id::UUID;
END;
$$;