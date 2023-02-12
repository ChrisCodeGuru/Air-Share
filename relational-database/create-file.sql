CREATE OR REPLACE FUNCTION create_file(
    _account_id VARCHAR,
    _name VARCHAR,
    _encrypted_key BYTEA,
    _asymmetric_key VARCHAR,
    _hash_value VARCHAR,
    _folder_id VARCHAR DEFAULT NULL
)
RETURNS UUID
LANGUAGE plpgsql
AS $$
DECLARE
    _new_id UUID;
    _parent_folder_permissions record;
BEGIN
    INSERT INTO file(folder_id, name, encrypted_key, asymmetric_key, hash_value)
    VALUES(_folder_id::UUID, _name, _encrypted_key, _asymmetric_key, _hash_value)
    RETURNING id INTO _new_id;

    IF _folder_id IS NULL THEN
        INSERT INTO account_file(account_id, file_id, permission)
        VALUES(_account_id::UUID, _new_id, 4);
    ELSE
        FOR _parent_folder_permissions IN
        (
            SELECT account_id, permission
            FROM account_folder
            WHERE folder_id = _folder_id::UUID
        )
        LOOP
            INSERT INTO account_file(account_id, file_id, permission)
            VALUES(_parent_folder_permissions.account_id, _new_id, _parent_folder_permissions.permission);
        END LOOP;
    END IF;

    RETURN _new_id;
END;
$$;