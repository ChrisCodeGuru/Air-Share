CREATE OR REPLACE FUNCTION create_folder(
    _account_id VARCHAR,
    _name VARCHAR,
    _folder_id VARCHAR DEFAULT NULL
)
RETURNS UUID
LANGUAGE plpgsql
AS $$
DECLARE
    _new_id UUID;
    _parent_folder_permissions record;
BEGIN
    -- Permissions to create in parent folder not checked yet
    -- We assume that the user is creating this folder in his own logical file storage system and not one that is shared with him
    INSERT INTO folder(name, parent_id)
    VALUES(_name, _folder_id::UUID)
    RETURNING id INTO _new_id;

    IF _folder_id IS NULL THEN
        INSERT INTO account_folder(account_id, folder_id, permission)
        VALUES(_account_id::UUID, _new_id, 4);
    ELSE
        FOR _parent_folder_permissions IN
        (
            SELECT account_id, permission
            FROM account_folder
            WHERE folder_id = _folder_id::UUID
        )
        LOOP
            INSERT INTO account_folder(account_id, folder_id, permission)
            VALUES(_parent_folder_permissions.account_id, _new_id, _parent_folder_permissions.permission);
        END LOOP;
    END IF;

    RETURN _new_id;
END;
$$;