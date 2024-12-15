



// get the time right now in America/New_York timezone
const getTimeInNewYork = () => {
  const newYorkTime = new Date().toLocaleString("en-US", { timeZone: "America/New_York" });
  return new Date(newYorkTime);
}

const schedules = [
  {
    "name": "business-hours",
    "startHour": 5,
    "endHour": 17,
    "days": [
      1,
      2,
      3,
      4,
      5
    ],
    "replicas": 3
  },
  {
    "name": "after-hours",
    "startHour": 17,
    "endHour": 5,
    "days": [
      1,
      2,
      3,
      4,
      5
    ],
    "replicas": 1
  }
]

const isTimeInSchedule = (schedule, time) => {
  const day = time.getDay();
  const hour = time.getHours();

  if (schedule.days.includes(day)) {
    return hour >= schedule.startHour && hour < schedule.endHour;
  }

  return false;
}
const loop = () => {
  console.log("Checking schedule...");
  const currentTime = getTimeInNewYork();
  for (const schedule of schedules) {
    if (isTimeInSchedule(schedule, currentTime)) {
      console.log(`[${currentTime}] Scaling to ${schedule.replicas} replicas`);
    }
  }
  setTimeout(loop, 1000);
}

loop();





