CREATE OR REPLACE FUNCTION delete_file(
    _file_id VARCHAR
)
RETURNS VARCHAR
LANGUAGE plpgsql
AS $$
DECLARE
    _deleted_name VARCHAR;
BEGIN
    DELETE FROM file 
    WHERE id = _file_id::UUID
    RETURNING name into _deleted_name;

    RETURN _deleted_name;
END;
$$;