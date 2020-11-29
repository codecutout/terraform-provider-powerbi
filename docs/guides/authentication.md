# Authentication
To use Power BI REST APIs the following need to be configured
1. A user with admin permission within Power BI. A user can be created by a domain owner, then Power BI admin permissions can be assigned via Office 365 admin center (link to the Office 365 admin center can be found within Power BI admin portal under users)
1. An Azure Active Directory App registration with delegate permissions on 'Power BI Service"' for 'Content.Create'. Easiest way to configure this is via https://dev.powerbi.com/apps 