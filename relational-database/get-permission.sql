CREATE OR REPLACE FUNCTION get_folder_permission(
    _folder_id VARCHAR
)
RETURNS TABLE (
    _account_id UUID,
    _account_email VARCHAR,
    _permission SMALLINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT account.id, account.email, account_folder.permission
    FROM account
    INNER JOIN account_folder
        ON account.id = account_folder.account_id
    WHERE account_folder.folder_id = _folder_id::UUID
    ORDER BY account_folder.permission DESC;
END;
$$;

CREATE OR REPLACE FUNCTION get_file_permission(
    _file_id VARCHAR
)
RETURNS TABLE (
    _account_id UUID,
    _account_email VARCHAR,
    _permission SMALLINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT account.id, account.email, account_file.permission
    FROM account
    INNER JOIN account_file
        ON account.id = account_file.account_id
    WHERE account_file.file_id = _file_id::UUID
    ORDER BY account_file.permission DESC;
END;
$$;