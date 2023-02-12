# Create directory for certs if it does not exist
$FolderName = ".\load-balancer\certs"
if (Test-Path $FolderName) {
    Write-Host ".\load-balancer\certs"
    # Perform Delete file from folder operation
}
else 
{
    # Powershell create directory if it does not exist
    New-Item $FolderName -ItemType Directory
    Write-Host "Folder created successfully"
}

# Generate SSL certs and Private key
Set-Location .\load-balancer\certs

# Specify the location of the installed openssl folder location. For example, ";C:\Users\en749\Downloads\openssl-0.9.8k_X64\bin"
$env:path = $env:path + ";C:\Users\en749\Downloads\openssl-0.9.8k_X64\bin"

# Specify the location of your openssl.cnf file. For example, "C:\Users\en749\Downloads\openssl-0.9.8k_X64\openssl.cnf"
$env:OPENSSL_CONF = "C:\Users\en749\Downloads\openssl-0.9.8k_X64\openssl.cnf"

# Generate localhost.key (key file) & localhost.csr file (request file)
openssl req -new -subj "/C=US/ST=Utah/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr

# Generate localhost.crt file (cert file)
openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt


# Add SSL cert to trusted cert list
# certutil -addstore -f "ROOT" app/certs/localhost.crt


#Reference Website:
# 1. https://stackoverflow.com/questions/63588254/how-to-set-up-an-https-server-with-a-self-signed-certificate-in-golang
# 2. https://stackoverflow.com/questions/14459078/unable-to-load-config-info-from-usr-local-ssl-openssl-cnf-on-windows
# 3. https://stackoverflow.com/questions/45506735/powershell-doesnt-recognize-openssl-even-after-i-added-it-to-the-system-path 