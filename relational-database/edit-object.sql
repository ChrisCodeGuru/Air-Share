CREATE OR REPLACE FUNCTION edit_folder(
    _folder_id VARCHAR,
    _folder_name VARCHAR
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE folder
    SET name = _folder_name
    WHERE id = _folder_id::UUID;

    RETURN true;
END;
$$;

CREATE OR REPLACE FUNCTION edit_file(
    _file_id VARCHAR,
    _file_name VARCHAR
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE file 
    SET name = _file_name
    WHERE id = _file_id::UUID;

    RETURN true;
END;
$$;