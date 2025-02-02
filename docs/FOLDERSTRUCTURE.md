# ServerCommander - Folder Structure

This document explains the folder structure of **ServerCommander** to help contributors and developers understand the project's organization.

```bash
ServerCommander/
├── src/
│   ├── rsrc.syso
│   ├── main.go          
│   ├── cli/ 
│   │   ├── cmd/  
│   │   │   ├── cli.go    
│   │   │   ├── server.go      
│   │   │   ├── process.go    
│   │   │   ├── file.go         
│   │   │   └── utils.go        
│   │   └── cli.go           
│   ├── assets/               
│   │   └── icon.ico
│   ├── services/  
│   │   └── logger.go  
│   ├── ui/                 
│   │   ├── applicationBanner.go 
│   │   ├── goodbyeBanner.go           
│   │   ├── upcomingFeaturesBanner.go                
│   │   └── welcomeCLIBanner.go                           
│   └── utils/             
│       ├── colors.go         
│       └── fileExists.go    
├── scripts/          
│   └── build.bat
├── docs/         
│   ├── README.md        
│   ├── INSTALL.md       
│   ├── USAGE.md        
│   ├── API.md           
│   ├── CONFIGURATION.md   
│   ├── THEMES.md     
│   ├── CONTRIBUTING.md  
│   ├── FOLDERSTRUCTURE.md 
│   └── CHANGELOG.md  
├── .gitignore   
├── go.mod    
├── LICENSE   
└── README.md 
```

## Explanation of Key Directories

**src/**
    - The **main directory** of the project, containing the entire source code and key files.

**cmd/**
    - Contains **command logic** and input handling for the CLI. Files like ```clear.go```, ```exit.go```, and ```help.go``` manage user interaction.

**assets/**
    - Stores **assets** such as ```icon.ico```, used for the user interface or the application.

**services/**
    - Provides core **service functions**, like logging (```logger.go```) and other important system services.

**ui/**
    - Contains **UI components** for the console and other banners like ```applicationBanner.go```, ```goodbyeBanner.go```, ```and welcomeCLIBanner.go```.

**utils/**
    - Provides **utility functions**, like color definitions (```colors.go```) and file existence checks (```fileExists.go```), to ease the program’s operation.

**scripts/**
    - Contains **automation scripts** for tasks like building the project (e.g., build.bat).

**docs/**
    - Stores documentation related to installation, usage, API, and other relevant information for developers and contributors.

## How to Navigate the Project

- If you want to **modify** the core application, check ```cmd/``` and ```services/```.
- If you want to **extend SSH or FTP functionality**, look into ```handlers/``` (this structure doesn’t exist yet in your current setup, but should be added for future expansion).
- If you want to **change configuration or themes**, check ```config/``` (also a future addition).
- If you want to **contribute tests**, extend the ```tests/``` folder (this folder is also not present yet).
- If you want to **understand the folder structure**, refer to ```docs/FOLDERSTRUCTURE.md``` (this document).

## Contribution

If you want to contribute, check out the guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).
