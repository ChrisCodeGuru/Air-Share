CREATE OR REPLACE FUNCTION update_password(
    _password VARCHAR,
    _account_id VARCHAR
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE account 
    SET password = _password
    WHERE id = _account_id::UUID;

    RETURN true;
END;
$$;