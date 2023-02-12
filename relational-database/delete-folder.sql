CREATE OR REPLACE FUNCTION delete_folder(
    _folder_id VARCHAR
)
RETURNS VARCHAR
LANGUAGE plpgsql
AS $$
DECLARE
    _deleted_name VARCHAR;
BEGIN
    DELETE FROM folder 
    WHERE id = _folder_id::UUID
    RETURNING name into _deleted_name;

    RETURN _deleted_name;
END;
$$;