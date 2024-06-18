# go version 1.21.5

# Development

## Build and Run 

 `go build -o bin/dev-app ./cmd/dev`

## Check Running Processes 

- Create file or edit new service file
`sudo nano /etc/systemd/system/devapp.service`

- content 
```
[Unit]
Description=Go DevApp Service
After=network.target

[Service]
Type=simple
ExecStart=/home/samsul-dev/projects/budget_in_backend/bin/devapp
WorkingDirectory=/home/samsul-dev/projects/budget_in_backend
EnvironmentFile=/home/samsul-dev/projects/budget_in_backend/dev.env
Restart=on-failure
User=samsul-dev
Group=samsul-dev

[Install]
WantedBy=multi-user.target
```

- Reload and Start the Service
```
sudo systemctl daemon-reload
sudo systemctl start devapp.service
sudo systemctl enable devapp.service
```

## Killing the Process 

- Check status
`sudo systemctl status devapp.service`

- Stop service
`sudo systemctl stop devapp.service`

- Disable the Service
`sudo systemctl disable devapp.service`

# Delete the Service

- Remove the Service file
`sudo rm /etc/systemd/system/devapp.service`

- Verify the Service is Removed
`systemctl status devapp.service`

# Production

## Build and Run 
1. `go build -o bin/prod-app ./cmd/prod`

## Check Running Processes 

- Create file or edit new service file
`sudo nano /etc/systemd/system/prodapp.service`

- content 
```
[Unit]
Description=Go ProdApp Service
After=network.target

[Service]
Type=simple
ExecStart=/home/samsul-dev/projects/budget_in_backend/bin/prod-app
WorkingDirectory=/home/samsul-dev/projects/budget_in_backend
EnvironmentFile=/home/samsul-dev/projects/budget_in_backend/prod.env
Restart=on-failure
User=samsul-dev
Group=samsul-dev

[Install]
WantedBy=multi-user.target
```

- Reload and Start the Service
```
sudo systemctl daemon-reload
sudo systemctl start prodapp.service
sudo systemctl enable prodapp.service
```

## Killing the Process 

- Check status
`sudo systemctl status prodapp.service`

- Stop service
`sudo systemctl stop prodapp.service`

- Disable the Service
`sudo systemctl disable prodapp.service`

# Delete the Service

- Remove the Service file
`sudo rm /etc/systemd/system/prodapp.service`

- Verify the Service is Removed
`systemctl status prodapp.service`