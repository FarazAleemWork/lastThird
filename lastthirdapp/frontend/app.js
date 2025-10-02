document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('locForm');
  const result = document.getElementById('result');

  if (!form || !result) {
    console.error('Missing form or result element');
    return;
  }

  const submitBtn = form.querySelector('button[type="submit"]');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    result.textContent = '';
    if (submitBtn) submitBtn.disabled = true;

    const cityVal = (document.getElementById('city') || {}).value?.trim() || '';
    const stateVal = (document.getElementById('state') || {}).value?.trim() || '';
    const countryVal = (document.getElementById('country') || {}).value?.trim() || '';
    const tzEl = document.getElementById('timezone');
    const tzVal = tzEl ? tzEl.value : '';

    if (!cityVal || !stateVal || !countryVal || !tzVal) {
      result.textContent = 'Please enter city, state, country and select a timezone.';
      if (submitBtn) submitBtn.disabled = false;
      return;
    }

    result.textContent = 'Loading...';

    const city = encodeURIComponent(cityVal);
    const state = encodeURIComponent(stateVal);
    const country = encodeURIComponent(countryVal);
    const timezone = encodeURIComponent(tzVal);

    const endpoint = `/api/geocode?city=${city}&state=${state}&country=${country}&timezone=${timezone}`;

    try {
      const res = await fetch(endpoint, { method: 'GET' });
      if (!res.ok) {
        const txt = await res.text().catch(() => '');
        result.textContent = `Server error: ${res.status} ${txt}`;
        return;
      }

      const data = await res.json().catch(() => null);
      if (!data) {
        result.textContent = 'Invalid JSON response from server';
        return;
      }

      const entry = Array.isArray(data) ? data[0] : data;
      const startTime = entry && (
        entry['Tahajjud starts at'] ??
        entry.tahajjudStart ??
        entry.tahajjud ??
        entry.start ??
        entry.time
      );

      if (startTime == null) {
        result.textContent = 'Response did not include Tahajjud start time';
        return;
      }

      result.textContent = `Tahajjud starts at: ${startTime}`;
    } catch (err) {
      result.textContent = `Request failed: ${err?.message || err}`;
    } finally {
      if (submitBtn) submitBtn.disabled = false;
    }
  });
});