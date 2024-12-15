const log = (level, message, extra) => {
  formatted = {
    severity: level.toUpperCase(),
    message,
    timestamp: new Date().toISOString(),
    ...(extra || {}),
  }
  console.log(JSON.stringify(formatted));
}

log("info", "booting");
let schedules = []
try {
  schedules = JSON.parse(process.env.SCHEDULES)
} catch (e) {
  log("error", "failed to parse schedules, did you set the SCHEDULES env var? it must be a valid json array", { error: e.message });
  process.exit(1);
}
let namespace = process.env.NAMESPACE;
if (!namespace) {
  log("error", "no namespace provided, did you set the NAMESPACE env var?");
  process.exit(1);
}
let deployment = process.env.DEPLOYMENT;
if (!deployment) {
  log("error", "no deployment provided, did you set the DEPLOYMENT env var?");
  process.exit(1);
}

log("info", "booted, config is as follows", { schedules, namespace, deployment });

const getTimeInNewYork = () => {
  const newYorkTime = new Date().toLocaleString("en-US", { timeZone: "America/New_York" });
  return new Date(newYorkTime);
}

const isTimeInSchedule = (schedule, time) => {
  const day = time.getDay();
  const hour = time.getHours();

  if (schedule.days.includes(day)) {
    return hour >= schedule.startHour && hour < schedule.endHour;
  }

  return false;
}

const loop = () => {
  log("debug", "HEARTBEAT, checking schedules");
  const currentTime = getTimeInNewYork();
  for (const schedule of schedules) {
    if (isTimeInSchedule(schedule, currentTime)) {
      log("info", `Matched schedule '${schedule.name}', scaling ${namespace}/${deployment} to ${schedule.replicas} replicas`, { schedule });
    }
  }
  setTimeout(loop, 1000);
}

loop();





