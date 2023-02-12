CREATE OR REPLACE FUNCTION check_existing_user(             -- Stored procedure to check if username or email is taken
    _email VARCHAR,
    _username VARCHAR
)
RETURNS TABLE (
    _count_email BIGINT,
    _count_username BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
    COUNT(email) AS _count_email,
    COUNT(username) AS _count_username
    FROM account
    WHERE email = _email
    OR username = _username;
END;
$$;