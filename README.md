# Teamdrive Manager (with extra Features)

## Needs
- GSuite Account (With a ton of permissions)

### Google Setup
- Go to the Dev Console of Google (https://console.developers.google.com/)
- Create a new API Project
    - Name it as you want in this tutorial I name it "TeamdriveManager"
    - After its created select it
- Click on "Enable APIs"
    - Enable the `Google Drive API`
    - Enable the `Admin SDK`
    - Enable the `Identity and Access Management (IAM) API`
    - Enable the `Cloud Resource Manager API`
    - Enable the `Service Management API`
    - Enable the `IAM Service Account Credentials API`
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
    - copy the Client ID to notepad.exe or so
- Go to the Admin Console (admin.google.com/YOURDOMAIN)
    - Open menu on the left
    - Go into "Security > Settings" (or use the search bar)
    - Click on the "Advanced Settings"
    - Select "Manage domain wide delegation" in the "Domain wide delegation" section
    - Click "Add new"
    - In the "Client ID" field enter the service account’s "Client ID"
    - In the next field, "OAuth scopes (comma-delimited)", enter the following 
    - `https://www.googleapis.com/auth/admin.directory.group,https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/drive,https://www.googleapis.com/auth/service.management`
