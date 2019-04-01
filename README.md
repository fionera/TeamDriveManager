# Teamdrive Manager (with extra Features)

## Needs
- GSuite Account (With a ton of permissions)

### Google Setup
- Go to the Dev Console of Google (https://console.developers.google.com/)
- Create a new API Project
    - Name it as you want in this tutorial I name it "TeamdriveManager"
    - After its created select it
- Click on "Enable APIs"
    - Enable the Admin SDK
    - Enable the Identity and Access Management (IAM) API
    - Enable the Google Drive API
- Click on "Credentials"
    - "Create Credentials"
    - "Service Account Key"
    - Create a new Service Account
    - As name you should use "TeamdriveManager-Impersonate"
    - Dont select a Role
    - As Type select JSON
    - When asked say "Create without Role"
    - You will now download a JSON File. DONT LOSE THE JSON FILE!
- Click on "Manage Service Accounts"
    - click on the mail address of the Service Account
    - Click Edit in the Top
    - Click on "Show Domain-wide delegation"
    - Enable "Enable G Suite Domain-wide Delegation"
    - As Product name just use the Project name again
    - Press Save
    - copy the Client ID to some notepad.exe or so
- Go to the Admin Console (admin.google.com/YOURDOMAIN)
    - Go into "Security" (or use the search bar)
    - Select "Show more" and then "Advanced settings"
    - Select "Manage API client access" in the "Authentication" section
    - In the "Client Name" field enter the service account’s "Client ID"
    - In the next field, "One or More API Scopes", enter the following 
    - `https://www.googleapis.com/auth/admin.directory.group,https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/drive`
