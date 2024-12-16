#!/bin/bash

# Schedule configurations JSON array
scheduleConfigsJson='[
  {
    "name": "weekday-evening",
    "startHour": 18,
    "endHour": 22,
    "days": [1, 2, 3, 4, 5],
    "replicas": 1
  },
  {
    "name": "weekend",
    "startHour": 0,
    "endHour": 24,
    "days": [6, 7],
    "replicas": 2
  }
]'

echo "Booting..."
echo "Raw Schedule configurations:"
echo "$scheduleConfigsJson"

echo "Parsed Schedule configurations:"
echo "$scheduleConfigsJson" | jq . 2>&1
if [ $? -ne 0 ]; then
  echo "Invalid JSON"
  exit 1
fi
echo "Boot Successful"

while true; do
  # Get current day and hour
  currentDay=$(TZ="America/New_York" date +%u)  # Day of the week (1-7, Monday is 1)
  currentHour=$(TZ="America/New_York" date +%H) # Hour of the day (00-23)

  echo "----------------------------------"
  echo "Current day: $currentDay"
  echo "Current hour: $currentHour"

  # Iterate over each scheduleConfig
  echo "$scheduleConfigsJson" | jq -c '.[]' | while read -r scheduleConfigJson; do
    echo "----------------------------------"
    # Extract fields from the current scheduleConfig
    scheduleName=$(echo "$scheduleConfigJson" | jq -r '.name')
    startHour=$(echo "$scheduleConfigJson" | jq -r '.startHour')
    endHour=$(echo "$scheduleConfigJson" | jq -r '.endHour')
    days=$(echo "$scheduleConfigJson" | jq -r '.days')
    replicas=$(echo "$scheduleConfigJson" | jq -r '.replicas')

    echo "Checking schedule: $scheduleName"
    echo "Start hour: $startHour, End hour: $endHour"
    echo "Days: $days, Replicas: $replicas"

    if [[ $currentHour -ge $startHour && $currentHour -lt $endHour && $days == *$currentDay* ]]; then
      echo "Scaling for $scheduleName to $replicas replicas"
      echo kubectl scale deployment my-deployment --replicas=$replicas
      break
    else
      echo "Not scaling for $scheduleName"
    fi
  done
  sleep 10
done
