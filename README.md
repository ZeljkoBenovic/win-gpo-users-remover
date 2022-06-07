# Windows GPO local users remover
This program removes all local users from a workstation, disables built-in Administrator user account
and creates a new user that is added to the Administrators group - local administrator user.    
It is intended to be used with GPO policy that runs this program on computer startup.   
In this way, a domain network has a unified local administrator account across all workstations
in case that a connection to the AD becomes unavailable for example.   

Existing user folders/profiles will not be deleted or changed in any way, only user accounts will be removed.

## How to use
Domain administrator can just place the `local-users-remove.exe` in this repo, to the location from 
which GPO will execute this program.   
Policy that runs scripts on startup can be found at `Computer Configuration > Policies > Windows Settings > Scripts (Startup/Shutdown) > Startup`    
Without defined parameters the program will create an administrator user using the default parameters `User: SecureAdmin / Pass: P@ssw0rd`   
In order to customize username and password, the program needs to be executed with the following parameters:   
`--admin-name <USERNAME>` - username for new administrator account    
`--admin-pass <PASSWORD>` - password for new administrator account
