const FormatDuration = (runtime) => {
  if (!runtime || runtime <= 0) {
    return undefined;
  }

  const hours = Math.floor(runtime / 60);
  const minutes = runtime % 60;

  if (hours > 0 && minutes > 0) {
    return `PT${hours}H${minutes}M`;
  }

  if (hours > 0) {
    return `PT${hours}H`;
  }

  return `PT${minutes}M`;
};

export default FormatDuration;