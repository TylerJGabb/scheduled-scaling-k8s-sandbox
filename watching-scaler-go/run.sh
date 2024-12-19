export SCHEDULES=$(cat example-schedules.json) 
export LOCAL=true 
export NAMESPACE=cronjob-sandbox 
export DEPLOYMENT=the-name-from-values-file-deploy 
go run .