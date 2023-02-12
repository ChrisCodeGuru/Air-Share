CREATE OR REPLACE FUNCTION sensitive_update(
    _file_id VARCHAR
)
RETURNS VARCHAR
LANGUAGE plpgsql
AS $$
DECLARE
    _file_name VARCHAR;
BEGIN
    UPDATE file 
    SET sensitive = true
    WHERE id = _file_id::UUID
    RETURNING name into _file_name;

    RETURN _file_name;
END;
$$;