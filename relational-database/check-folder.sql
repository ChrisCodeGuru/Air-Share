CREATE OR REPLACE FUNCTION check_folder_existence(
    _folder_id VARCHAR
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    _count BIGINT;
BEGIN
    SELECT COUNT(*) 
    INTO _count
    FROM folder 
    WHERE id = _folder_id::UUID;

    RETURN _count;
END;
$$;